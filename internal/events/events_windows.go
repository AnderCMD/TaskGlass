//go:build windows

// Package events turns Windows window events and a global hotkey into cheap
// channel signals, so TaskGlass reacts instantly without polling in a loop.
//
// A dedicated OS thread runs a message loop that owns the WinEvent hooks and
// the hotkey registration (both are thread-affine).
package events

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procSetWinEventHook   = user32.NewProc("SetWinEventHook")
	procUnhookWinEvent    = user32.NewProc("UnhookWinEvent")
	procGetMessageW       = user32.NewProc("GetMessageW")
	procPeekMessageW      = user32.NewProc("PeekMessageW")
	procPostThreadMessage = user32.NewProc("PostThreadMessageW")
	procRegisterHotKey    = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey  = user32.NewProc("UnregisterHotKey")
	procGetCurrentThread  = kernel32.NewProc("GetCurrentThreadId")
)

const (
	wmHotkey = 0x0312
	wmQuit   = 0x0012
	wmApp    = 0x8000 // wParam = modifiers, lParam = virtual key (0 = unregister)

	modNoRepeat = 0x4000
	hotkeyID    = 1

	winEventOutOfContext = 0x0000
	winEventSkipOwn      = 0x0002

	eventSystemForeground = 0x0003
	eventSystemMinStart   = 0x0016
	eventSystemMinEnd     = 0x0017
	eventObjectDestroy    = 0x8001
	eventObjectLocation   = 0x800B
)

type point struct{ X, Y int32 }
type msg struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

// Watcher delivers window-change and hotkey signals.
type Watcher struct {
	// Changed receives a value (coalesced) whenever the foreground window,
	// window visibility, minimise state or window position changes.
	Changed chan struct{}
	// Hotkey receives a value whenever the registered hotkey is pressed.
	Hotkey chan struct{}

	log    *slog.Logger
	thread uint32
	done   chan struct{}
}

// package-level because a WinEvent callback has no user-data pointer.
var current *Watcher

var winEventCB = syscall.NewCallback(func(hook, event, hwnd, idObject, idChild, thread, ts uintptr) uintptr {
	w := current
	if w == nil {
		return 0
	}
	// Object events fire for every control and caret; only whole windows matter.
	if event >= eventObjectDestroy && (int32(idObject) != 0 || int32(idChild) != 0) {
		return 0
	}
	select {
	case w.Changed <- struct{}{}:
	default: // already pending
	}
	return 0
})

// Start launches the event thread. Only one Watcher may exist at a time.
func Start(log *slog.Logger) (*Watcher, error) {
	if current != nil {
		return nil, errors.New("events: already started")
	}
	w := &Watcher{
		Changed: make(chan struct{}, 1),
		Hotkey:  make(chan struct{}, 1),
		log:     log,
		done:    make(chan struct{}),
	}
	current = w
	ready := make(chan error, 1)
	go w.run(ready)
	select {
	case err := <-ready:
		if err != nil {
			current = nil
			return nil, err
		}
	case <-time.After(5 * time.Second):
		current = nil
		return nil, errors.New("events: thread did not start")
	}
	return w, nil
}

func (w *Watcher) run(ready chan<- error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(w.done)

	// Touch the message queue so PostThreadMessage has somewhere to deliver.
	var m msg
	procPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, 0)
	tid, _, _ := procGetCurrentThread.Call()
	w.thread = uint32(tid)

	var hooks []uintptr
	for _, r := range [][2]uintptr{
		{eventSystemForeground, eventSystemForeground},
		{eventSystemMinStart, eventSystemMinEnd},
		{eventObjectDestroy, eventObjectLocation},
	} {
		h, _, err := procSetWinEventHook.Call(r[0], r[1], 0, winEventCB, 0, 0, winEventOutOfContext|winEventSkipOwn)
		if h == 0 {
			ready <- fmt.Errorf("SetWinEventHook(%#x..%#x): %w", r[0], r[1], err)
			return
		}
		hooks = append(hooks, h)
	}
	defer func() {
		for _, h := range hooks {
			procUnhookWinEvent.Call(h)
		}
		procUnregisterHotKey.Call(0, hotkeyID)
	}()
	ready <- nil

	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 { // WM_QUIT or error
			return
		}
		switch m.Message {
		case wmHotkey:
			select {
			case w.Hotkey <- struct{}{}:
			default:
			}
		case wmApp:
			procUnregisterHotKey.Call(0, hotkeyID)
			if m.LParam != 0 {
				ok, _, err := procRegisterHotKey.Call(0, hotkeyID, m.WParam|modNoRepeat, m.LParam)
				if ok == 0 {
					w.log.Warn("hotkey could not be registered (is it used by another program?)", "err", err)
				}
			}
		}
	}
}

// SetHotkey registers (or, with vk == 0, removes) the global toggle hotkey.
// mods and vk use Win32 values; see config.ParseHotkey.
func (w *Watcher) SetHotkey(mods, vk uint32) {
	procPostThreadMessage.Call(uintptr(w.thread), wmApp, uintptr(mods), uintptr(vk))
}

// Stop ends the event thread and removes every hook.
func (w *Watcher) Stop() {
	procPostThreadMessage.Call(uintptr(w.thread), wmQuit, 0, 0)
	select {
	case <-w.done:
	case <-time.After(2 * time.Second):
	}
	current = nil
}
