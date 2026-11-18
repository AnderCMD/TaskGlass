//go:build windows

// Package xaml removes the Windows 11 taskbar's own XAML background fill and
// can make the tray overflow flyout translucent.
//
// The modern taskbar paints a semi-transparent grey fill over whatever the DWM
// accent draws, so accent effects alone can never reach "fully transparent",
// and the flyout with the hidden tray icons draws its own opaque background.
// To change them, TaskGlass loads a small native DLL (native/tap) into
// explorer.exe through the XAML Diagnostics API. This package gets that DLL
// loaded with a WH_CALLWNDPROC hook on the taskbar's thread (the DLL then
// attaches itself to XAML from inside Explorer) and tells it, through named
// shared memory, what to do.
//
// If the DLL is missing or injection fails, TaskGlass keeps working with
// accent effects only.
package xaml

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

// DLLName is the TAP library shipped next to taskglass.exe.
const DLLName = "taskglass_tap.dll"

const (
	stateName = `Local\TaskGlass.Xaml.State.v1`
	readyName = `Local\TaskGlass.Xaml.Ready`
)

// Slots of the shared block (8 int32s). Keep in sync with native/tap/tap.c.
// Slots 3 to 6 are unused and must stay 0: copies of the helper DLL from
// unreleased builds, still loaded in a running Explorer or Start host, read them
// as "surface opacity" and a non-zero value would change those surfaces.
const (
	slotHide     = 0 // TaskGlass -> TAP: 1 hides the taskbar fill
	slotStatus   = 1 // TAP -> TaskGlass: 1 attached, 2 controlling, <0 failed
	slotOverflow = 2 // TaskGlass -> TAP: tray overflow flyout, see Encode
	slotFound    = 7 // TAP -> TaskGlass: 1 once the tray overflow flyout is tracked
	numSlots     = 8
)

// Encode converts a surface setting to the shared-memory encoding: 0 leaves
// the surface alone, otherwise 1 + the opacity percent (0..100).
func Encode(enabled bool, opacityPercent int) int32 {
	if !enabled {
		return 0
	}
	if opacityPercent < 0 {
		opacityPercent = 0
	}
	if opacityPercent > 100 {
		opacityPercent = 100
	}
	return int32(1 + opacityPercent)
}

// injectTimeout bounds how long we wait for the DLL to attach to Explorer
// (it retries for up to 30 s while XAML starts up).
const injectTimeout = 35 * time.Second

// Controller injects the TAP and drives what it does.
type Controller struct {
	log     *slog.Logger
	dllPath string

	vals  *[numSlots]int32 // shared memory, nil until Start succeeds
	view  uintptr
	hmap  windows.Handle
	ready windows.Handle

	hide atomic.Bool

	cancel context.CancelFunc
	done   chan struct{}

	pid       uint32 // explorer pid we last attached to
	failLog   string
	probes    int // tray probes tried for this pid
	lastProbe time.Time
	probedPid uint32
}

// Available reports whether the TAP DLL exists next to the executable and
// returns its path.
func Available(exe string) (string, bool) {
	p := filepath.Join(filepath.Dir(exe), DLLName)
	_, err := os.Stat(p)
	return p, err == nil
}

// New creates a Controller for the DLL at dllPath.
func New(log *slog.Logger, dllPath string) *Controller {
	return &Controller{log: log, dllPath: dllPath}
}

func (c *Controller) load(slot int) int32     { return atomic.LoadInt32(&c.vals[slot]) }
func (c *Controller) store(slot int, v int32) { atomic.StoreInt32(&c.vals[slot], v) }

// Start creates the shared block and starts the injector loop.
func (c *Controller) Start() error {
	name, _ := windows.UTF16PtrFromString(stateName)
	size := uint32(numSlots * 4)
	h, err := windows.CreateFileMapping(windows.InvalidHandle, nil, windows.PAGE_READWRITE, 0, size, name)
	if err != nil {
		return fmt.Errorf("create shared state: %w", err)
	}
	view, err := windows.MapViewOfFile(h, windows.FILE_MAP_WRITE, 0, 0, uintptr(size))
	if err != nil {
		windows.CloseHandle(h)
		return fmt.Errorf("map shared state: %w", err)
	}
	c.hmap, c.view = h, view
	c.vals = (*[numSlots]int32)(*(*unsafe.Pointer)(unsafe.Pointer(&view)))
	for i := range c.vals {
		c.store(i, 0)
	}
	evName, _ := windows.UTF16PtrFromString(readyName)
	ev, err := windows.CreateEvent(nil, 1, 0, evName) // manual reset
	if err != nil {
		c.release()
		return fmt.Errorf("create ready event: %w", err)
	}
	c.ready = ev

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.done = make(chan struct{})
	go c.run(ctx)
	return nil
}

// SetHide asks the TAP to hide (true) or restore (false) the taskbar's XAML
// fill. Injection only happens the first time something needs it.
func (c *Controller) SetHide(hide bool) {
	if c.vals == nil {
		return
	}
	c.hide.Store(hide)
	v := int32(0)
	if hide {
		v = 1
	}
	c.store(slotHide, v)
}

// SetOverflow sets the tray overflow flyout (see Encode).
func (c *Controller) SetOverflow(v int32) {
	if c.vals == nil {
		return
	}
	c.store(slotOverflow, v)
}

// Stop restores everything, gives the TAP time to notice, and releases
// resources.
func (c *Controller) Stop() {
	if c.vals == nil {
		return
	}
	c.SetHide(false)
	c.SetOverflow(0)
	c.cancel()
	<-c.done
	time.Sleep(500 * time.Millisecond) // TAP polls every 150 ms, then dispatches to the UI thread
	c.release()
}

func (c *Controller) release() {
	if c.vals == nil {
		return
	}
	windows.UnmapViewOfFile(c.view)
	windows.CloseHandle(c.hmap) // once gone, the TAP also falls back to "leave alone"
	if c.ready != 0 {
		windows.CloseHandle(c.ready)
		c.ready = 0
	}
	c.vals = nil
}

func (c *Controller) needed() bool {
	return c.hide.Load() || c.load(slotOverflow) != 0
}

func (c *Controller) run(ctx context.Context) {
	defer close(c.done)
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		c.ensureInjected(ctx)
		c.ensureOverflow(ctx)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (c *Controller) ensureInjected(ctx context.Context) {
	if !c.needed() {
		return
	}
	bars := winapi.Taskbars()
	if len(bars) == 0 {
		return
	}
	tid, pid := winapi.WindowThreadPID(bars[0])
	if tid == 0 || pid == 0 || pid == c.pid {
		return
	}
	path, err := stage(c.dllPath)
	if err == nil {
		err = c.inject(ctx, bars[0], tid, pid, path)
	}
	if err != nil {
		if msg := err.Error(); msg != c.failLog {
			c.failLog = msg
			c.log.Warn("xaml surface control unavailable (will retry)", "err", err)
		}
		return
	}
	c.pid, c.failLog = pid, ""
	c.log.Info("xaml tap attached", "explorer_pid", pid)
	if c.hide.Load() {
		c.awaitFills(ctx)
	}
}

// awaitFills waits until the TAP reports it controls the taskbar fill (status
// 2). The TAP can only see taskbar elements through live tree events, so if
// none happens we provoke one with a short-lived window.
func (c *Controller) awaitFills(ctx context.Context) {
	for attempt := 0; attempt < 6; attempt++ {
		for i := 0; i < 6; i++ { // 1.5 s
			if c.load(slotStatus) == 2 {
				c.log.Info("xaml taskbar fill is being controlled")
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(250 * time.Millisecond):
			}
		}
		if err := winapi.TaskbarAnchor(300 * time.Millisecond); err != nil {
			c.log.Debug("anchor window failed", "err", err)
		}
	}
	c.log.Warn("xaml taskbar fill not found yet; it will be picked up on the next taskbar change")
}

// ensureOverflow makes the helper DLL find the tray overflow flyout. The flyout
// is a XAML island that exists from the first time it is opened; if it was
// created before the DLL attached, the DLL receives no event for it. Adding and
// removing a tray icon makes the tray emit events, which lets the DLL find it.
// It tries a few times per Explorer process and then gives up (the flyout is
// then picked up the next time Explorer restarts or its tray changes).
func (c *Controller) ensureOverflow(ctx context.Context) {
	if c.load(slotOverflow) == 0 || c.load(slotFound) == 1 || c.pid == 0 {
		return
	}
	if c.probedPid != c.pid {
		c.probedPid, c.probes = c.pid, 0
	}
	if c.probes >= 4 || time.Since(c.lastProbe) < 6*time.Second {
		return
	}
	c.probes++
	c.lastProbe = time.Now()
	if err := winapi.TrayProbe(400 * time.Millisecond); err != nil {
		c.log.Debug("tray probe failed", "err", err)
		return
	}
	c.log.Debug("tray probe sent to wake the overflow flyout", "attempt", c.probes)
	select {
	case <-ctx.Done():
	case <-time.After(1500 * time.Millisecond):
	}
	if c.load(slotFound) == 1 {
		c.log.Info("tray overflow flyout is being controlled")
	}
}

// stage copies the DLL to a content-addressed path under the user cache so
// that an updated DLL never collides with a copy Explorer still has loaded.
func stage(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := hex.EncodeToString(h.Sum(nil))[:12]
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "TaskGlass", "tap")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, "taskglass_tap_"+sum+".dll")
	if _, err := os.Stat(dst); err == nil {
		return dst, nil
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(out, f); err != nil {
		out.Close()
		os.Remove(tmp)
		return "", err
	}
	if err := out.Close(); err != nil {
		return "", err
	}
	return dst, os.Rename(tmp, dst)
}

// inject gets the TAP DLL loaded into the taskbar's thread of Explorer and
// waits for it to report that it attached to XAML.
func (c *Controller) inject(ctx context.Context, bar winapi.HWND, tid, pid uint32, dllPath string) error {
	if moduleLoaded(pid, filepath.Base(dllPath)) {
		// An earlier TaskGlass run attached this exact DLL, or a failed attempt
		// is still unwinding. Only the TAP can tell: it republishes status 2 to
		// every new shared block while it controls the taskbar fill.
		for i := 0; i < 12; i++ { // 3 s
			if c.load(slotStatus) == 2 {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(250 * time.Millisecond):
			}
		}
		return fmt.Errorf("tap is loaded in Explorer but not controlling the taskbar yet")
	}
	windows.ResetEvent(c.ready)
	c.store(slotStatus, 0)

	// Load the DLL here too: the hook needs a module handle and the address of
	// the hook procedure. Its DllMain does nothing outside explorer.exe.
	mod, err := windows.LoadLibrary(dllPath)
	if err != nil {
		return fmt.Errorf("load tap: %w", err)
	}
	defer windows.FreeLibrary(mod)
	proc, err := windows.GetProcAddress(mod, "TaskGlassHookProc")
	if err != nil {
		return fmt.Errorf("tap has no hook procedure: %w", err)
	}
	hook, err := winapi.SetCallWndProcHook(proc, mod, tid)
	if err != nil {
		return fmt.Errorf("install hook: %w", err)
	}
	defer winapi.Unhook(hook)
	winapi.NotifyNull(bar) // make the hook fire so Windows loads the DLL in Explorer

	deadline := time.Now().Add(injectTimeout)
	for time.Now().Before(deadline) {
		ev, _ := windows.WaitForSingleObject(c.ready, 250)
		if ev == windows.WAIT_OBJECT_0 {
			if st := c.load(slotStatus); st != 1 && st != 2 {
				return fmt.Errorf("tap failed to attach: status 0x%08X (see tap.log)", uint32(st))
			}
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return fmt.Errorf("tap did not attach within %s (see tap.log)", injectTimeout)
}

// moduleLoaded reports whether process pid has a module with the given file name.
func moduleLoaded(pid uint32, name string) bool {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE|windows.TH32CS_SNAPMODULE32, pid)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(snap)
	var me windows.ModuleEntry32
	me.Size = uint32(unsafe.Sizeof(me))
	for err = windows.Module32First(snap, &me); err == nil; err = windows.Module32Next(snap, &me) {
		if strings.EqualFold(windows.UTF16ToString(me.Module[:]), name) {
			return true
		}
	}
	return false
}
