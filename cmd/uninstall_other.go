//go:build !windows

package cmd

import "syscall"

// detachedProcAttr is unused off Windows (Unix can unlink a running binary),
// but removeSelf references it, so it must exist for every platform.
func detachedProcAttr() *syscall.SysProcAttr { return nil }
