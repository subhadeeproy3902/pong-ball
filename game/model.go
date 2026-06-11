package game

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/subhadeeproy3902/paddle-ball/store"
	"github.com/subhadeeproy3902/paddle-ball/ui"
)

// BaseSpeed is the ball speed in cells/second at Phase 1.
const BaseSpeed = 18.0

// ────────────────────────────────────────────────────────────────────────────
// Model
// ────────────────────────────────────────────────────────────────────────────

// Model is the top-level bubbletea model; it holds all application state.
type Model struct {
	// ── window ──
	width, height int
	playW, playH  int // play-area dimensions (inside borders + header/footer)

	// ── screen state machine ──
	phase   AppPhase
	mode    GameMode
	menuSel int // title screen cursor (0-3)

	// ── game objects ──
	ball      Ball
	paddle    Paddle
	particles []Particle
	ftexts    []FloatText
	fallingPUs []FallingPU
	activePU  *ActivePU

	// ── input state ──
	keyUp, keyDown bool

	// ── scoring ──
	score     int
	hiScore   int
	streak    int
	maxStreak int
	catches   int
	misses    int

	// ── difficulty ──
	curPhase Phase

	// ── lives (Arcade) ──
	lives int

	// ── timing ──
	lastTick  time.Time
	gameStart time.Time
	elapsed   time.Duration
	timeLimit time.Duration // Time Trial only

	// ── countdown ──
	countdown int     // 3, 2, 1
	cdTTL     float64 // seconds until next decrement

	// ── phase-transition banner ──
	bannerText  string
	bannerColor string
	bannerTTL   float64

	// ── iron shield state ──
	shieldActive bool

	// ── ghost state ──
	ghostActive bool

	// ── theme ──
	themeIdx int

	// ── power-up spawn counter ──
	catchesSinceLastPU int

	// ── store / leaderboard ──
	st       *store.Store
	scores   []store.ScoreRecord
	lbFilter string // "" = all, else mode code
}

// ────────────────────────────────────────────────────────────────────────────
// Constructor
// ────────────────────────────────────────────────────────────────────────────

// NewModel creates a fresh model, pre-loading the hi-score and config.
func NewModel(modeStr, themeStr string) Model {
	st := store.New()

	// resolve theme
	themeIdx := ui.ThemeIndexByName(themeStr)

	// resolve mode (may be "" for title-screen flow)
	mode := ModeClassic
	startOnTitle := true
	if modeStr != "" {
		mode = parseModeStr(modeStr)
		startOnTitle = false
	}

	hi := st.HiScore(modeStr)

	m := Model{
		phase:    PhaseTitle,
		mode:     mode,
		menuSel:  0,
		themeIdx: themeIdx,
		st:       st,
		hiScore:  hi,
	}

	if !startOnTitle {
		m.startCountdown()
	}

	return m
}

func parseModeStr(s string) GameMode {
	switch strings.ToLower(s) {
	case "arcade":
		return ModeArcade
	case "zen":
		return ModeZen
	case "timed", "time-trial", "timetrial":
		return ModeTimeTrial
	}
	return ModeClassic
}

// ────────────────────────────────────────────────────────────────────────────
// bubbletea interface
// ────────────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return schedTick()
}

// ────────────────────────────────────────────────────────────────────────────
// Internal helpers — game initialization
// ────────────────────────────────────────────────────────────────────────────

// startCountdown transitions to the 3-2-1-GO phase.
func (m *Model) startCountdown() {
	m.phase = PhaseCountdown
	m.countdown = 3
	m.cdTTL = 1.0
	m.resetGameObjects()
}

// startPlaying begins the actual game.
func (m *Model) startPlaying() {
	m.phase = PhasePlaying
	m.gameStart = time.Now()
	m.lastTick = time.Now()
	m.elapsed = 0
	m.score = 0
	m.streak = 0
	m.maxStreak = 0
	m.catches = 0
	m.misses = 0
	m.particles = nil
	m.ftexts = nil
	m.fallingPUs = nil
	m.activePU = nil
	m.shieldActive = false
	m.ghostActive = false
	m.catchesSinceLastPU = 0
	m.bannerTTL = 0
	m.curPhase = Phases[0]

	switch m.mode {
	case ModeArcade:
		m.lives = 3
	case ModeTimeTrial:
		m.lives = 999
		m.timeLimit = 60 * time.Second
	default:
		m.lives = 1
	}

	m.hiScore = m.st.HiScore(m.mode.ShortCode())
	m.resetGameObjects()
}

// resetGameObjects re-spawns the ball and paddle at neutral positions.
func (m *Model) resetGameObjects() {
	if m.playH == 0 {
		m.playH = 20 // sensible fallback before first resize
		m.playW = 70
	}

	ph := PhaseForScore(m.score)
	m.curPhase = ph

	paddleH := ph.PaddleH
	// keep existing paddle H if we're mid-game and just lost a life
	if m.paddle.H > 0 && m.phase == PhasePlaying {
		paddleH = m.paddle.H
	}

	m.paddle = Paddle{
		Y: float64(m.playH/2 - paddleH/2),
		H: paddleH,
	}

	// ball starts in the middle, heading right at a slight downward angle
	speed := BaseSpeed * ph.SpeedMult
	angle := -0.3 + rand.Float64()*0.6 // small random angle
	m.ball = Ball{
		X:  float64(m.playW/2),
		Y:  float64(m.playH / 2),
		VX: speed * math.Cos(angle),
		VY: speed * math.Sin(angle),
	}
	// ensure ball moves rightward first (away from paddle) to give player time
	if m.ball.VX < 0 {
		m.ball.VX = -m.ball.VX
	}
}

// ────────────────────────────────────────────────────────────────────────────
// schedTick
// ────────────────────────────────────────────────────────────────────────────

func schedTick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}