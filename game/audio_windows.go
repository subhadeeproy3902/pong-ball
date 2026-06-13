//go:build windows

package game

// Windows audio backend — MCI (mciSendStringW from winmm.dll), CGO-free.
//
// WHY A DEDICATED THREAD: the "mpegvideo" MCI device is DirectShow/COM-based. It
// creates a hidden window on the thread that OPENS it, and that thread must
// stay COM-initialised (STA) and PUMP window messages for playback to actually
// run. Go schedules goroutines across OS threads freely, so the old design
// (open on one goroutine, play from transient goroutines, nobody pumping)
// silently produced no sound at all. Here every MCI call happens on ONE locked
// OS thread that pumps messages; play requests arrive over a channel.
//
// VOICE POOL: frequent sounds (paddle hit, wall bounce) get several pre-cued
// copies, round-robined, so a rapid rally never waits on a single device's
// rewind — zero-lag overlapping playback.

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

var (
	winmm          = syscall.NewLazyDLL("winmm.dll")
	procMCISendStr = winmm.NewProc("mciSendStringW")

	ole32      = syscall.NewLazyDLL("ole32.dll")
	procCoInit = ole32.NewProc("CoInitializeEx")

	user32           = syscall.NewLazyDLL("user32.dll")
	procPeekMessage  = user32.NewProc("PeekMessageW")
	procTranslateMsg = user32.NewProc("TranslateMessage")
	procDispatchMsg  = user32.NewProc("DispatchMessageW")

	playReq = make(chan Sfx, 64)
)

// winMSG mirrors the Win32 MSG struct (x64 layout; releases build windows/amd64).
type winMSG struct {
	hwnd     uintptr
	message  uint32
	wParam   uintptr
	lParam   uintptr
	time     uint32
	ptX      int32
	ptY      int32
	lPrivate uint32
}

func mci(command string) {
	ptr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return
	}
	// mciSendStringW(command, returnBuf=NULL, returnLen=0, callbackHwnd=NULL)
	procMCISendStr.Call(uintptr(unsafe.Pointer(ptr)), 0, 0, 0)
}

// pumpMessages drains the thread's window-message queue so the MCI device's
// hidden window keeps servicing playback.
func pumpMessages() {
	var msg winMSG
	for {
		r, _, _ := procPeekMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 1) // PM_REMOVE
		if r == 0 {
			return
		}
		procTranslateMsg.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMsg.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// backendInit launches the dedicated audio worker, which opens the devices.
func backendInit(paths map[Sfx]string) {
	go audioWorker(paths)
}

func audioWorker(paths map[Sfx]string) {
	runtime.LockOSThread() // pin to one OS thread for the lifetime of the worker
	procCoInit.Call(0, 0x2) // CoInitializeEx(NULL, COINIT_APARTMENTTHREADED); ignore result

	type pool struct {
		aliases []string
		next    int
	}
	pools := make(map[Sfx]*pool)

	idx := 0
	for sfx, p := range paths {
		voices := 1
		if sfx == SfxHit || sfx == SfxBounce {
			voices = 4 // frequent → pool so overlapping replays never wait on a rewind
		}
		pl := &pool{}
		for v := 0; v < voices; v++ {
			alias := fmt.Sprintf("pbsfx%d", idx)
			idx++
			mci("close " + alias) // drop any stale alias from a previous run
			mci(fmt.Sprintf(`open "%s" type mpegvideo alias %s`, p, alias))
			mci("cue " + alias + " output") // prime: no warm-up lag on first play
			pl.aliases = append(pl.aliases, alias)
		}
		pools[sfx] = pl
	}

	// Pump messages continuously (so the device windows stay live) and play
	// requests as they arrive on the channel.
	ticker := time.NewTicker(4 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case s := <-playReq:
			if pl := pools[s]; pl != nil && len(pl.aliases) > 0 {
				alias := pl.aliases[pl.next]
				pl.next = (pl.next + 1) % len(pl.aliases)
				mci("play " + alias + " from 0") // async; restarts from the top
			}
		case <-ticker.C:
		}
		pumpMessages()
	}
}

// backendPlay hands the event to the audio worker without blocking the render
// loop. If the worker is still opening devices (or flooded), the request is
// dropped rather than stalling the game.
func backendPlay(s Sfx, _ string) {
	select {
	case playReq <- s:
	default:
	}
}
