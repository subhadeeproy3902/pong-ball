package game

// sound.go — event → sound mapping and throttling. The actual playback lives in
// audio.go (+ the platform backends). Sounds are real MP3s from soundcn, not the
// terminal bell.

// Sfx identifies a game sound event.
type Sfx int

const (
	SfxHit    Sfx = iota // paddle catch
	SfxMiss              // ball lost
	SfxPower             // power-up collected
	SfxPhase             // difficulty phase up
	SfxStart             // countdown GO
	SfxOver              // game over
	SfxBest              // new personal best
	SfxMenu              // menu move / selection
	SfxBounce            // ball bounces off a wall
)

// requestSfx plays a sound for an event, honouring the mute toggle. SfxHit is
// rate-limited so a fast rally is rhythmic rather than a continuous blast.
// Playback is dispatched to a goroutine so it never stalls the render loop.
func (m *Model) requestSfx(s Sfx) {
	if !m.soundOn {
		return
	}
	switch s {
	case SfxHit:
		if m.hitBellCD > 0 {
			return
		}
		m.hitBellCD = HitBellGap
	case SfxBounce:
		if m.bounceCD > 0 {
			return
		}
		m.bounceCD = BounceGap
	}
	go playSfx(s)
}
