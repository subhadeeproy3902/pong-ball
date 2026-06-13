//go:build windows

package game

// Windows audio backend — MCI (Media Control Interface) through winmm.dll.
// Pure syscall, no CGO. Each sound opens once as an aliased device; play
// retriggers from the start and is asynchronous, and distinct sounds overlap.

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

var (
	winmm          = syscall.NewLazyDLL("winmm.dll")
	procMCISendStr = winmm.NewProc("mciSendStringW")

	mciMu    sync.Mutex
	mciAlias = map[Sfx]string{}
)

func mci(command string) {
	ptr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return
	}
	// mciSendStringW(command, returnBuf=NULL, returnLen=0, callbackHwnd=NULL)
	procMCISendStr.Call(uintptr(unsafe.Pointer(ptr)), 0, 0, 0)
}

func backendInit(paths map[Sfx]string) {
	mciMu.Lock()
	defer mciMu.Unlock()
	i := 0
	for sfx, p := range paths {
		alias := fmt.Sprintf("pbsfx%d", i)
		i++
		mci("close " + alias) // drop any stale alias from a previous run
		mci(fmt.Sprintf(`open "%s" type mpegvideo alias %s`, p, alias))
		mciAlias[sfx] = alias
	}
}

func backendPlay(s Sfx, _ string) {
	mciMu.Lock()
	alias, ok := mciAlias[s]
	mciMu.Unlock()
	if !ok {
		return
	}
	mci("play " + alias + " from 0") // async; restarts if already playing
}
