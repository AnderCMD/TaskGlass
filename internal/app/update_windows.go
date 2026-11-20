//go:build windows

package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/settings"
	"github.com/AnderCMD/TaskGlass/internal/update"
)

const (
	firstUpdateCheck = 90 * time.Second // let start-up settle first
	updateInterval   = 24 * time.Hour
)

// updater holds the update state of the running app.
type updater struct {
	mu         sync.Mutex
	client     *update.Client
	latest     update.Release
	haveLatest bool
	installing bool
	asked      string // version the user was already asked about (notify mode)
	failed     string // version whose install failed; not retried until restart
}

func (a *App) updateClient() *update.Client {
	if a.upd.client == nil {
		a.upd.client = update.New("TaskGlass/" + a.Version)
	}
	return a.upd.client
}

// CheckUpdate asks GitHub for the latest release and remembers it, so a later
// InstallUpdate installs exactly what was reported.
func (a *App) CheckUpdate() (update.Info, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.upd.mu.Lock()
	c := a.updateClient()
	a.upd.mu.Unlock()
	rel, err := c.Latest(ctx)
	if err != nil {
		return update.Info{Current: a.Version}, err
	}
	info := update.Info{Current: a.Version, Latest: rel.Version, Notes: rel.Notes, Available: update.Newer(a.Version, rel.Version)}
	a.upd.mu.Lock()
	if info.Available {
		a.upd.latest, a.upd.haveLatest = rel, true
	} else {
		a.upd.haveLatest = false
	}
	a.upd.mu.Unlock()
	a.Log.Info("update check", "current", a.Version, "latest", rel.Version, "available", info.Available)
	return info, nil
}

// InstallUpdate downloads the release found by the last CheckUpdate, verifies
// it and starts its installer, then quits so the installer can replace the
// program. If the installer cannot be started (for example blocked by security
// software) the verified file is kept and its path is copied to the clipboard.
func (a *App) InstallUpdate() error {
	a.upd.mu.Lock()
	if a.upd.installing {
		a.upd.mu.Unlock()
		return errors.New("an update is already being installed")
	}
	if !a.upd.haveLatest {
		a.upd.mu.Unlock()
		return errors.New("no update has been found yet")
	}
	rel, c := a.upd.latest, a.updateClient()
	a.upd.installing = true
	a.upd.mu.Unlock()
	done := func(failedVersion string) {
		a.upd.mu.Lock()
		a.upd.installing = false
		if failedVersion != "" {
			a.upd.failed = failedVersion
		}
		a.upd.mu.Unlock()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	path, err := c.Download(ctx, rel, filepath.Join(os.TempDir(), "TaskGlass-update"))
	if err != nil {
		done(rel.Version)
		a.Log.Error("update download", "version", rel.Version, "err", err)
		return err
	}
	if err := update.RunInstaller(path); err != nil {
		done(rel.Version)
		a.Log.Error("update installer could not be started", "path", path, "err", err)
		settings.CopyToClipboard(path)
		return fmt.Errorf(a.tr("update.manual"), path)
	}
	a.Log.Info("update installer started; quitting", "version", rel.Version)
	go func() {
		time.Sleep(700 * time.Millisecond)
		a.Quit()
	}()
	return nil
}

// updateLoop checks for updates shortly after start and then every day, as
// configured. In "notify" mode it asks once per version; in "auto" mode it
// installs without asking.
func (a *App) updateLoop(ctx context.Context) {
	t := time.NewTimer(firstUpdateCheck)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
		t.Reset(updateInterval)
		a.backgroundUpdateCheck()
	}
}

func (a *App) backgroundUpdateCheck() {
	mode := a.snapshot().Updates
	if mode == config.UpdatesOff {
		return
	}
	info, err := a.CheckUpdate()
	if err != nil {
		a.Log.Warn("update check failed", "err", err)
		return
	}
	if !info.Available {
		return
	}
	a.upd.mu.Lock()
	skip := a.upd.failed == info.Latest || a.upd.installing
	ask := mode == config.UpdatesNotify && a.upd.asked != info.Latest
	if ask {
		a.upd.asked = info.Latest
	}
	a.upd.mu.Unlock()
	if skip {
		return
	}
	switch mode {
	case config.UpdatesAuto:
		if err := a.InstallUpdate(); err != nil {
			a.Log.Error("automatic update failed", "err", err)
			a.warnText(fmt.Sprintf(a.tr("update.failed"), err.Error()))
		}
	case config.UpdatesNotify:
		if ask {
			a.askInstall(info)
		}
	}
}

// askInstall shows the "install now?" question and installs on yes.
func (a *App) askInstall(info update.Info) {
	if !settings.Ask("TaskGlass", fmt.Sprintf(a.tr("update.available"), info.Latest, info.Current)) {
		return
	}
	if err := a.InstallUpdate(); err != nil {
		a.warnText(fmt.Sprintf(a.tr("update.failed"), err.Error()))
	}
}

// manualUpdateCheck is the tray item "Check for updates…": it always answers.
func (a *App) manualUpdateCheck() {
	info, err := a.CheckUpdate()
	switch {
	case err != nil:
		a.warnText(fmt.Sprintf(a.tr("update.checkFailed"), err.Error()))
	case info.Available:
		a.askInstall(info)
	default:
		a.warnText(fmt.Sprintf(a.tr("update.uptodate"), info.Current))
	}
}

// warnText shows a message box without blocking the caller.
func (a *App) warnText(text string) {
	go settings.Notify("TaskGlass", strings.TrimSpace(text))
}
