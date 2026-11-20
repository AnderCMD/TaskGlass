//go:build windows

package settings

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// edgePath finds Microsoft Edge (shipped with Windows 10/11).
func edgePath() string {
	for _, root := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		k, err := registry.OpenKey(root, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe`, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		v, _, err := k.GetStringValue("")
		k.Close()
		if err == nil {
			if _, err := os.Stat(v); err == nil {
				return v
			}
		}
	}
	for _, env := range []string{"ProgramFiles(x86)", "ProgramFiles"} {
		if base := os.Getenv(env); base != "" {
			p := filepath.Join(base, `Microsoft\Edge\Application\msedge.exe`)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

// Show opens url in a standalone app-style window (no tabs or address bar).
// It uses Edge with a private profile so it neither joins the user's browser
// session nor shows first-run prompts; if Edge is missing it falls back to the
// default browser.
func Show(url, profileDir string) error {
	if edge := edgePath(); edge != "" {
		cmd := exec.Command(edge,
			"--app="+url,
			"--user-data-dir="+profileDir,
			"--window-size=1240,820",
			"--no-first-run",
			"--no-default-browser-check",
			"--disable-extensions",
		)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		if err := cmd.Start(); err == nil {
			go cmd.Wait()
			return nil
		}
	}
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url)
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}
