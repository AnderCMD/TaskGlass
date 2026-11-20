//go:build windows

package update

import (
	"os/exec"
	"syscall"
)

const (
	detachedProcess       = 0x00000008
	createNewProcessGroup = 0x00000200
)

// RunInstaller starts the downloaded installer silently and returns without
// waiting: the installer closes the running TaskGlass itself (and starts the new
// one when it is done), so the caller should exit right after.
func RunInstaller(path string) error {
	cmd := exec.Command(path, "/SILENT", "/SUPPRESSMSGBOXES", "/NORESTART", "/CLOSEAPPLICATIONS")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: detachedProcess | createNewProcessGroup}
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
