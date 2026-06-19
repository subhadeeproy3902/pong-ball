package cmd

import "syscall"

// Windows constants: run the cleanup helper without a console window and
// detached from this process so it survives our exit.
const (
	detachedProcess = 0x00000008
	createNoWindow  = 0x08000000
)

func detachedProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: detachedProcess | createNoWindow}
}
