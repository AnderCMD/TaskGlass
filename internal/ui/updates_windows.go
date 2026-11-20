//go:build windows

package ui

import "fmt"

// checkUpdates asks the host for the latest release without blocking the UI.
func (w *window) checkUpdates() {
	w.updMsg.SetText(w.t("update.checking"))
	w.updCheck.SetEnabled(false)
	go func() {
		info, err := w.host.CheckUpdate()
		w.mw.Synchronize(func() {
			w.updCheck.SetEnabled(true)
			switch {
			case err != nil:
				w.updMsg.SetText(fmt.Sprintf(w.t("update.checkFailed"), err.Error()))
			case info.Available:
				w.updMsg.SetText(fmt.Sprintf(w.t("update.availableShort"), info.Latest))
				w.updInstall.SetVisible(true)
			default:
				w.updMsg.SetText(fmt.Sprintf(w.t("update.uptodate"), info.Current))
				w.updInstall.SetVisible(false)
			}
		})
	}()
}

// installUpdate downloads, verifies and starts the installer found by the last check.
func (w *window) installUpdate() {
	w.updMsg.SetText(w.t("update.installing"))
	w.updInstall.SetEnabled(false)
	go func() {
		err := w.host.InstallUpdate()
		w.mw.Synchronize(func() {
			w.updInstall.SetEnabled(true)
			if err != nil {
				w.updMsg.SetText(fmt.Sprintf(w.t("update.failed"), err.Error()))
			}
		})
	}()
}
