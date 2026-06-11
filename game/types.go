package game

import "time"

// ────────────────────────────────────────────────────────────────────────────
// Enums / type aliases
// ────────────────────────────────────────────────────────────────────────────

// TickMsg drives the 60 fps game loop.
type TickMsg time.Time

// AppPhase is the current screen shown to the user.
type AppPhase int

const (
	PhaseTitle       AppPhase = iota // main menu
	PhaseCountdown                   // 3 … 2 … 1 … GO!
	PhasePlaying                     // active gameplay
	PhasePaused                      // paused overlay
	PhaseGameOver                    // end-of-round breakdown
	PhaseLeaderboard                 // score history table
	PhaseHelp                        // controls reference
)

// GameMode selects which rules to use.
type GameMode int

const (
	ModeClassic  GameMode = iota // one life, pure score chase
	ModeArcade                   // 3 lives + power-ups
	ModeZen                      // infinite lives, just vibe
	ModeTimeTrial                // 60-second race
)

func (m GameMode) String() string {
	switch m {
	case ModeClassic:
		return "Classic"
	case ModeArcade:
		return "Arcade"
	case ModeZen:
		return "Zen"
	case ModeTimeTrial:
		return "Time Trial"
	}
	return "Unknown"
}

func (m GameMode) ShortCode() string {
	switch m {
	case ModeClassic:
		return "classic"
	case ModeArcade:
		return "arcade"
	case ModeZen:
		return "zen"
	case ModeTimeTrial:
		return "timed"
	}
	return "classic"
}

// ────────────────────────────────────────────────────────────────────────────
// Game object types
// ────────────────────────────────────────────────────────────────────────────

// Pt is a 2-D integer coordinate.
type Pt struct{ X, Y int }

// Ball holds position, velocity, and trail history.
type Ball struct {
	X, Y   float64 // current position (float for sub-cell precision)
	VX, VY float64 // velocity in cells/second
	Trail  []Pt    // last N positions for rendering the heat trail
}

// Paddle holds state for the player paddle.
type Paddle struct {
	Y        float64 // top edge (float for smooth movement)
	H        int     // current height in cells
	VY       float64 // current vertical velocity (used for spin transfer)
	FlashTTL float64 // seconds of white flash remaining after a hit
}

// PaddleX is the fixed horizontal position of the paddle (column index in play area).
const PaddleX = 1

// ────────────────────────────────────────────────────────────────────────────
// Particle system
// ────────────────────────────────────────────────────────────────────────────

// Particle is one transient visual element (wall impact, score pop spread, etc.)
type Particle struct {
	X, Y   float64
	VX, VY float64
	Life   float64 // starts at 1.0, dies at 0
	Decay  float64 // life units lost per second
	Glyph  rune
	Color  string
}

// FloatText is a "+N" score label that drifts upward and fades.
type FloatText struct {
	X, Y  float64
	Text  string
	Color string
	Life  float64
	Decay float64
}

// ────────────────────────────────────────────────────────────────────────────
// Power-up system
// ────────────────────────────────────────────────────────────────────────────

// PowerUpKind identifies a power-up type.
type PowerUpKind int

const (
	PUWidePaddle PowerUpKind = iota // paddle grows 3 cells
	PUSlowMo                        // ball speed −35 %
	PUFirePaddle                    // score ×2 on every hit
	PUIronShield                    // one auto-save bounce off left wall
	PUGhost                         // pass-through once without penalty
	PUBomb                          // NEGATIVE — shrinks paddle −2 cells
)

func (k PowerUpKind) Glyph() rune {
	return []rune{'Ⓦ', 'ⓢ', 'ⓕ', 'ⓘ', 'ⓖ', 'Ⓑ'}[k]
}

func (k PowerUpKind) Name() string {
	return []string{"Wide Paddle", "Slow Mo", "Fire Paddle", "Iron Shield", "Ghost Ball", "BOMB"}[k]
}

// Duration returns effect duration in seconds; 0 = one-shot.
func (k PowerUpKind) Duration() float64 {
	return []float64{12, 8, 15, 0, 0, 10}[k]
}

func (k PowerUpKind) Color() string {
	return []string{"#89DDFF", "#C3E88D", "#FFD700", "#4ECDC4", "#F8F8F8", "#FF5370"}[k]
}

// FallingPU is a power-up that drops from the top of the play area.
type FallingPU struct {
	X, Y  float64
	Kind  PowerUpKind
	VY    float64 // falling speed cells/sec
	Frame int     // animation frame index
	FrameTTL float64
}

// ActivePU is a power-up currently in effect.
type ActivePU struct {
	Kind  PowerUpKind
	TTL   float64 // seconds remaining (−∞ = already consumed one-shot)
	Total float64 // original duration (for progress bar)
}

// ────────────────────────────────────────────────────────────────────────────
// Difficulty phase table
// ────────────────────────────────────────────────────────────────────────────

// Phase holds parameters for one difficulty level.
type Phase struct {
	Num       int
	Name      string
	Color     string
	Emoji     string
	MinScore  int
	SpeedMult float64 // multiplier on BaseSpeed
	PaddleH   int     // paddle height in cells
	TrailLen  int     // trail positions to render
}

// Phases lists all five difficulty levels in order.
var Phases = []Phase{
	{1, "Warm Up", "#C3E88D", "🌱", 0, 1.00, 7, 1},
	{2, "Heating Up", "#FFCB6B", "🔥", 10, 1.25, 6, 2},
	{3, "On Fire", "#FF8C00", "💥", 25, 1.55, 5, 3},
	{4, "Blazing", "#FF5370", "⚡", 50, 1.90, 4, 4},
	{5, "INSANE", "#FF00FF", "🏆", 100, 2.35, 3, 5},
}

// PhaseForScore returns the Phase that corresponds to the given score.
func PhaseForScore(score int) Phase {
	p := Phases[0]
	for _, ph := range Phases {
		if score >= ph.MinScore {
			p = ph
		}
	}
	return p
}