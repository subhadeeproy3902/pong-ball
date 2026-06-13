package game

// audio.go — real sound effects, sourced from soundcn (https://soundcn.xyz).
//
// The MP3s are embedded in the binary and extracted to a temp dir at first run.
// Playback is platform-specific and CGO-free:
//   - Windows: MCI via winmm.dll (mciSendStringW) — plays MP3 directly, overlaps
//     across distinct sounds, no external player, no process spawn.
//   - macOS / Linux: a detected CLI player (afplay / mpg123 / ffplay …), async.
// This avoids any audio C library so the cross-platform release stays CGO-free.

import (
	"embed"
	"os"
	"path/filepath"
	"sync"
)

//go:embed sounds/*.mp3
var soundFS embed.FS

// soundFile maps each event to its embedded soundcn MP3.
var soundFile = map[Sfx]string{
	SfxHit:    "coin-collect.mp3",         // paddle catch (short, frequent)
	SfxMiss:   "phaser-down-2.mp3",        // ball lost (descending)
	SfxPower:  "power-up-3.mp3",           // power-up collected
	SfxPhase:  "phase-jump-1.mp3",         // difficulty phase up
	SfxStart:  "begin.mp3",                // countdown GO
	SfxOver:   "explosion-crunch-002.mp3", // game over
	SfxBest:   "arcade-mode.mp3",          // new personal best (fanfare)
	SfxMenu:   "select-003.mp3",           // menu move / selection
	SfxBounce: "select-003.mp3",           // ball off a wall (reuse the crisp blip)
}

var (
	audioOnce  sync.Once
	audioReady bool
	sfxPath    = map[Sfx]string{}
)

// initAudio extracts the embedded sounds to a temp dir and prepares the backend.
// Runs exactly once. Callers should invoke it off the hot path (e.g. a goroutine
// at startup) since the platform backend may pre-open the files.
func initAudio() {
	audioOnce.Do(func() {
		dir := filepath.Join(os.TempDir(), "paddle-ball-sfx")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return
		}
		for sfx, name := range soundFile {
			data, err := soundFS.ReadFile("sounds/" + name)
			if err != nil {
				continue
			}
			p := filepath.Join(dir, name)
			if fi, statErr := os.Stat(p); statErr != nil || fi.Size() != int64(len(data)) {
				if os.WriteFile(p, data, 0o644) != nil {
					continue
				}
			}
			sfxPath[sfx] = p
		}
		backendInit(sfxPath)
		audioReady = true
	})
}

// playSfx plays the sound for an event (no-op until audio is ready).
func playSfx(s Sfx) {
	if !audioReady {
		return
	}
	if p, ok := sfxPath[s]; ok {
		backendPlay(s, p)
	}
}
