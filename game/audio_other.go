//go:build !windows

package game

// Non-Windows audio backend — shell out to a detected CLI player, async.
// macOS ships afplay (mp3); most Linux desktops have mpg123 or ffplay. If none
// is found, sound silently no-ops (the game is fully playable without it).

import "os/exec"

var playerCmd []string

func backendInit(_ map[Sfx]string) {
	for _, c := range [][]string{
		{"afplay"},       // macOS
		{"mpg123", "-q"}, // linux (mp3)
		{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet"}, // ffmpeg
		{"play", "-q"}, // sox
		{"cvlc", "--play-and-exit", "--intf", "dummy"}, // vlc
	} {
		if _, err := exec.LookPath(c[0]); err == nil {
			playerCmd = c
			return
		}
	}
}

func backendPlay(_ Sfx, path string) {
	if playerCmd == nil {
		return
	}
	args := append(append([]string{}, playerCmd[1:]...), path)
	cmd := exec.Command(playerCmd[0], args...)
	if cmd.Start() == nil {
		go func() { _ = cmd.Wait() }() // reap, don't block
	}
}
