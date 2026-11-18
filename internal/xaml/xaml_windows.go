//go:build windows

// Package xaml controls XAML surfaces of the Windows 11 shell that accent
// effects cannot reach.
//
// The modern taskbar paints a semi-transparent grey fill over whatever the DWM
// accent draws, and flyouts such as the tray overflow or the Start menu draw
// their own opaque backgrounds. To change them, TaskGlass loads a small native
// DLL (native/tap) into the process that owns each surface through the XAML
// Diagnostics API: explorer.exe (taskbar, tray overflow) and
// StartMenuExperienceHost.exe (Start menu). This package gets the DLL loaded
// with a WH_CALLWNDPROC hook (the DLL then attaches itself to XAML from inside
// the process) and tells it, through named shared memory, what to do.
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
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

// DLLName is the TAP library shipped next to taskglass.exe.
const DLLName = "taskglass_tap.dll"

const stateName = `Local\TaskGlass.Xaml.State.v1`

// Slots of the shared block (8 int32s). Keep in sync with native/tap/tap.c.
const (
	slotHide     = 0 // TaskGlass -> TAP: 1 hides the taskbar fill
	slotExplorer = 1 // TAP -> TaskGlass: explorer status (1 attached, 2 controlling, <0 failed)
	slotOverflow = 2 // TaskGlass -> TAP: tray overflow flyout, see Encode
	slotStart    = 3 // TaskGlass -> TAP: Start menu surfaces, see Encode
	slotShell    = 4 // TaskGlass -> TAP: notification centre, calendar, quick settings
	slotStartSt  = 5 // TAP -> TaskGlass: Start host status
	slotShellSt  = 6 // TAP -> TaskGlass: ShellExperienceHost status
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

// injectTimeout bounds how long we wait for the DLL to attach to a process
// (it retries for up to 30 s while XAML starts up).
const injectTimeout = 35 * time.Second

// target is a process whose XAML tree the TAP is loaded into.
type target struct {
	name         string
	readyName    string
	slot         int // status slot
	ready        windows.Handle
	find         func() (hwnd winapi.HWND, tid, pid uint32, ok bool)
	needed       func(c *Controller) bool
	hookProc     string // exported name of the hook procedure in the DLL
	getMsg       bool   // WH_GETMESSAGE on GUI threads instead of WH_CALLWNDPROC on a window
	threads      func() (pid uint32, tids []uint32, ok bool)
	appContainer bool // the process is an AppContainer: the DLL must live where it can read it

	pid     uint32 // pid we last attached to
	failLog string
}

// Controller injects the TAP and drives what it does.
type Controller struct {
	log     *slog.Logger
	dllPath string

	vals *[numSlots]int32 // shared memory, nil until Start succeeds
	view uintptr
	hmap windows.Handle

	targets []*target
	mu      sync.Mutex
	cancel  context.CancelFunc
	done    chan struct{}
	hide    atomic.Bool
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
	c := &Controller{log: log, dllPath: dllPath}
	c.targets = []*target{
		{
			name: "explorer", hookProc: "TaskGlassHookProc", readyName: `Local\TaskGlass.Xaml.Ready`, slot: slotExplorer,
			find: func() (winapi.HWND, uint32, uint32, bool) {
				bars := winapi.Taskbars()
				if len(bars) == 0 {
					return 0, 0, 0, false
				}
				tid, pid := winapi.WindowThreadPID(bars[0])
				return bars[0], tid, pid, tid != 0 && pid != 0
			},
			needed: func(c *Controller) bool {
				return c.hide.Load() || c.load(slotOverflow) != 0
			},
		},
		{
			name: "start", readyName: `Local\TaskGlass.Xaml.Ready.Start`, slot: slotStartSt,
			// The Start host owns no window while Start is closed, so hook its UI
			// thread directly.
			getMsg: true, hookProc: "TaskGlassGetMsgProc",
			threads: func() (uint32, []uint32, bool) { return winapi.ProcessGUIThreads("StartMenuExperienceHost.exe") },
			find: func() (winapi.HWND, uint32, uint32, bool) {
				pid, tids, ok := winapi.ProcessGUIThreads("StartMenuExperienceHost.exe")
				if !ok {
					return 0, 0, 0, false
				}
				return 0, tids[0], pid, true
			},
			needed: func(c *Controller) bool { return c.load(slotStart) != 0 },
		},
		{
			name: "shell", readyName: `Local\TaskGlass.Xaml.Ready.Shell`, slot: slotShellSt,
			getMsg: true, hookProc: "TaskGlassGetMsgProc", appContainer: true,
			threads: func() (uint32, []uint32, bool) { return winapi.ProcessGUIThreads("ShellExperienceHost.exe") },
			find: func() (winapi.HWND, uint32, uint32, bool) {
				pid, tids, ok := winapi.ProcessGUIThreads("ShellExperienceHost.exe")
				if !ok {
					return 0, 0, 0, false
				}
				return 0, tids[0], pid, true
			},
			needed: func(c *Controller) bool { return c.load(slotShell) != 0 },
		},
	}
	return c
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
	for _, t := range c.targets {
		n, _ := windows.UTF16PtrFromString(t.readyName)
		ev, err := windows.CreateEvent(nil, 1, 0, n) // manual reset
		if err != nil {
			c.release()
			return fmt.Errorf("create ready event: %w", err)
		}
		t.ready = ev
	}

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

// SetSurfaces sets the tray overflow, Start menu and shell panels (see Encode).
func (c *Controller) SetSurfaces(overflow, start, shell int32) {
	if c.vals == nil {
		return
	}
	c.store(slotOverflow, overflow)
	c.store(slotStart, start)
	c.store(slotShell, shell)
}

// Stop restores everything, gives the TAP time to notice, and releases
// resources.
func (c *Controller) Stop() {
	if c.vals == nil {
		return
	}
	c.SetHide(false)
	c.SetSurfaces(0, 0, 0)
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
	for _, t := range c.targets {
		if t.ready != 0 {
			windows.CloseHandle(t.ready)
			t.ready = 0
		}
	}
	c.vals = nil
}

func (c *Controller) run(ctx context.Context) {
	defer close(c.done)
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		for _, tg := range c.targets {
			c.ensureInjected(ctx, tg)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (c *Controller) ensureInjected(ctx context.Context, t *target) {
	if !t.needed(c) {
		return
	}
	hwnd, tid, pid, ok := t.find()
	if !ok {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if pid == t.pid {
		return
	}
	path, err := stage(c.dllPath, t.appContainer)
	if err == nil {
		err = c.inject(ctx, t, hwnd, tid, pid, path)
	}
	if err != nil {
		if msg := err.Error(); msg != t.failLog {
			t.failLog = msg
			c.log.Warn("xaml surface control unavailable (will retry)", "target", t.name, "err", err)
		}
		return
	}
	t.pid, t.failLog = pid, ""
	c.log.Info("xaml tap attached", "target", t.name, "pid", pid)
	if t.name == "explorer" && c.hide.Load() {
		c.awaitFills(ctx)
	}
}

// awaitFills waits until the TAP reports it controls the taskbar fill (status
// 2). The TAP can only see taskbar elements through live tree events, so if
// none happens we provoke one with a short-lived window.
func (c *Controller) awaitFills(ctx context.Context) {
	for attempt := 0; attempt < 6; attempt++ {
		for i := 0; i < 6; i++ { // 1.5 s
			if c.load(slotExplorer) == 2 {
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

// stage copies the DLL to a content-addressed path under the user cache so
// that an updated DLL never collides with a copy a process still has loaded.
func stage(src string, appContainer bool) (string, error) {
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
	if appContainer {
		// AppContainer processes cannot read the user profile; ProgramData is
		// readable once ALL APPLICATION PACKAGES is granted on the file.
		dir = filepath.Join(os.Getenv("ProgramData"), "TaskGlass", "tap")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, "taskglass_tap_"+sum+".dll")
	if _, err := os.Stat(dst); err == nil {
		return dst, nil
	}
	if appContainer {
		defer grantAppContainerRead(dir, dst)
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

// inject gets the TAP DLL loaded into the thread that owns hwnd and waits for
// it to report that it attached to XAML.
func (c *Controller) inject(ctx context.Context, t *target, hwnd winapi.HWND, tid, pid uint32, dllPath string) error {
	if moduleLoaded(pid, filepath.Base(dllPath)) {
		// An earlier TaskGlass run attached this exact DLL, or a failed attempt
		// is still unwinding. Only the TAP can tell: it republishes status 2 to
		// every new shared block while it controls something.
		for i := 0; i < 12; i++ { // 3 s
			if c.load(t.slot) == 2 {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(250 * time.Millisecond):
			}
		}
		return fmt.Errorf("tap is loaded in %s but not controlling anything yet", t.name)
	}
	windows.ResetEvent(t.ready)
	c.store(t.slot, 0)

	// Load the DLL here too: the hook needs a module handle and the address of
	// the hook procedure. Its DllMain does nothing outside the target processes.
	mod, err := windows.LoadLibrary(dllPath)
	if err != nil {
		return fmt.Errorf("load tap: %w", err)
	}
	defer windows.FreeLibrary(mod)
	proc, err := windows.GetProcAddress(mod, t.hookProc)
	if err != nil {
		return fmt.Errorf("tap has no hook procedure %s: %w", t.hookProc, err)
	}
	// The UI thread is not known in advance, so hook every thread that has a
	// message queue and wake each of them; the DLL is loaded once per process.
	var hooks []uintptr
	defer func() {
		for _, h := range hooks {
			winapi.Unhook(h)
		}
	}()
	if t.getMsg {
		_, tids, _ := t.threads()
		for _, id := range tids {
			if h, err := winapi.SetGetMessageHook(proc, mod, id); err == nil {
				hooks = append(hooks, h)
			}
		}
		if len(hooks) == 0 {
			return fmt.Errorf("install hook: no thread of %s accepted it", t.name)
		}
		for _, id := range tids {
			winapi.PostThreadNull(id)
		}
	} else {
		h, err := winapi.SetCallWndProcHook(proc, mod, tid)
		if err != nil {
			return fmt.Errorf("install hook: %w", err)
		}
		hooks = append(hooks, h)
		winapi.NotifyNull(hwnd) // make the hook fire so Windows loads the DLL there
	}

	deadline := time.Now().Add(injectTimeout)
	for time.Now().Before(deadline) {
		ev, _ := windows.WaitForSingleObject(t.ready, 250)
		if ev == windows.WAIT_OBJECT_0 {
			if st := c.load(t.slot); st != 1 && st != 2 {
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

// grantAppContainerRead lets AppContainer processes read the staged DLL: it
// grants "ALL APPLICATION PACKAGES" (S-1-15-2-1) read and execute on the file
// and on its folder.
func grantAppContainerRead(dir, file string) {
	for _, p := range []string{dir, file} {
		_ = exec.Command("icacls", p, "/grant", "*S-1-15-2-1:(RX)").Run()
	}
}
