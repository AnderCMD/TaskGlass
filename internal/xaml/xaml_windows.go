//go:build windows

// Package xaml removes the Windows 11 taskbar's own XAML background fill.
//
// The modern taskbar paints a semi-transparent grey fill over whatever the DWM
// accent draws, so accent effects alone can never reach "fully transparent".
// To remove it, TaskGlass loads a small native DLL (native/tap) into
// explorer.exe through the XAML Diagnostics API. This package gets that DLL
// loaded into Explorer's taskbar thread with a WH_CALLWNDPROC hook (the DLL then
// attaches itself to XAML from inside Explorer) and tells it, through a named
// shared-memory flag, whether to hide the fill.
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
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

// DLLName is the TAP library shipped next to taskglass.exe.
const DLLName = "taskglass_tap.dll"

const (
	stateName = `Local\TaskGlass.Xaml.State.v1` // two int32s: [0] mode, [1] status
	readyName = `Local\TaskGlass.Xaml.Ready`
)

// injectTimeout bounds how long we wait for the DLL to attach inside Explorer
// (it retries for up to 30 s while XAML starts up).
const injectTimeout = 35 * time.Second

// Controller injects the TAP and drives its hide/show flag.
type Controller struct {
	log     *slog.Logger
	dllPath string

	flag   *int32 // shared memory [0], nil until Start succeeds
	status *int32 // shared memory [1], written by the TAP
	view   uintptr
	hmap   windows.Handle
	ready  windows.Handle

	want    atomic.Bool
	mu      sync.Mutex
	pid     uint32 // explorer pid we last injected into
	failLog string
	cancel  context.CancelFunc
	done    chan struct{}
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

// Start creates the shared flag and starts the injector loop.
func (c *Controller) Start() error {
	name, _ := windows.UTF16PtrFromString(stateName)
	h, err := windows.CreateFileMapping(windows.InvalidHandle, nil, windows.PAGE_READWRITE, 0, 8, name)
	if err != nil {
		return fmt.Errorf("create shared state: %w", err)
	}
	view, err := windows.MapViewOfFile(h, windows.FILE_MAP_WRITE, 0, 0, 8)
	if err != nil {
		windows.CloseHandle(h)
		return fmt.Errorf("map shared state: %w", err)
	}
	c.hmap, c.view = h, view
	base := *(*unsafe.Pointer)(unsafe.Pointer(&view))
	c.flag = (*int32)(base)
	c.status = (*int32)(unsafe.Add(base, 4))
	atomic.StoreInt32(c.flag, 0)
	atomic.StoreInt32(c.status, 0)
	evName, _ := windows.UTF16PtrFromString(readyName)
	ev, err := windows.CreateEvent(nil, 1, 0, evName) // manual reset
	if err != nil {
		windows.UnmapViewOfFile(view)
		windows.CloseHandle(h)
		return fmt.Errorf("create ready event: %w", err)
	}
	c.ready = ev

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.done = make(chan struct{})
	go c.run(ctx)
	return nil
}

// SetHide asks the TAP to hide (true) or restore (false) the XAML fill.
// Injection only happens the first time hiding is requested.
func (c *Controller) SetHide(hide bool) {
	if c.flag == nil {
		return
	}
	c.want.Store(hide)
	v := int32(0)
	if hide {
		v = 1
	}
	atomic.StoreInt32(c.flag, v)
}

// Stop restores the fill, gives the TAP time to notice, and releases resources.
func (c *Controller) Stop() {
	if c.flag == nil {
		return
	}
	c.SetHide(false)
	c.cancel()
	<-c.done
	time.Sleep(500 * time.Millisecond) // TAP polls every 150 ms, then dispatches to the UI thread
	windows.UnmapViewOfFile(c.view)
	windows.CloseHandle(c.hmap) // once gone, the TAP also falls back to "leave alone"
	windows.CloseHandle(c.ready)
	c.flag = nil
}

func (c *Controller) run(ctx context.Context) {
	defer close(c.done)
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		c.ensureInjected(ctx)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (c *Controller) ensureInjected(ctx context.Context) {
	if !c.want.Load() {
		return
	}
	bars := winapi.Taskbars()
	if len(bars) == 0 {
		return
	}
	tid, pid := winapi.WindowThreadPID(bars[0])
	c.mu.Lock()
	defer c.mu.Unlock()
	if pid == 0 || pid == c.pid {
		return
	}
	path, err := stage(c.dllPath)
	if err == nil {
		err = c.inject(ctx, bars[0], tid, pid, path)
	}
	if err != nil {
		if msg := err.Error(); msg != c.failLog {
			c.failLog = msg
			c.log.Warn("xaml fill removal unavailable (will retry)", "err", err)
		}
		return
	}
	c.pid, c.failLog = pid, ""
	c.log.Info("xaml tap attached", "explorer_pid", pid)
	c.awaitFills(ctx)
}

// awaitFills waits until the TAP reports it controls the taskbar fill (status
// 2). The TAP can only see taskbar elements through live tree events, so if
// none happens we provoke one with a short-lived window.
func (c *Controller) awaitFills(ctx context.Context) {
	for attempt := 0; attempt < 6; attempt++ {
		for i := 0; i < 6; i++ { // 1.5 s
			if atomic.LoadInt32(c.status) == 2 {
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

// inject gets the TAP DLL loaded into Explorer's taskbar thread and waits for
// it to report that it attached to XAML.
func (c *Controller) inject(ctx context.Context, bar winapi.HWND, tid, pid uint32, dllPath string) error {
	if moduleLoaded(pid, filepath.Base(dllPath)) {
		return nil // an earlier TaskGlass run already attached this exact DLL
	}
	if tid == 0 {
		return fmt.Errorf("taskbar window has no thread")
	}
	windows.ResetEvent(c.ready)
	atomic.StoreInt32(c.status, 0)

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
			if st := atomic.LoadInt32(c.status); st != 1 {
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
