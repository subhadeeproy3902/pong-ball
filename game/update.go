package game

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/subhadeeproy3902/paddle-ball/store"
)

const numThemes = 4

// Update is the bubbletea update function — the heart of the game.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcPlayArea()
		if m.phase == PhasePlaying {
			m.clampPaddle()
		}
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		m = m.handleKey(msg)
	case TickMsg:
		m = m.tick(time.Time(msg))
	}
	return m, schedTick()
}

// ─────────────────────────────────────────────────────────────────────────────
// Key handling
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleKey(k tea.KeyMsg) Model {
	switch m.phase {
	case PhaseTitle:
		m.keyUp = false
		m.keyDown = false
		switch k.String() {
		case "up", "k", "w":
			m.menuSel = (m.menuSel + 3) % 4
		case "down", "j", "s":
			m.menuSel = (m.menuSel + 1) % 4
		case "enter", " ":
			m.mode = GameMode(m.menuSel)
			m.startCountdown()
		case "1":
			m.mode = ModeClassic
			m.startCountdown()
		case "2":
			m.mode = ModeArcade
			m.startCountdown()
		case "3":
			m.mode = ModeZen
			m.startCountdown()
		case "4":
			m.mode = ModeTimeTrial
			m.startCountdown()
		case "S":
			m.phase = PhaseLeaderboard
			recs, _ := m.st.LoadAll()
			m.scores = recs
		case "t", "T":
			m.themeIdx = (m.themeIdx + 1) % numThemes
		case "?", "h":
			m.phase = PhaseHelp
		case "q":
			return m
		}

	case PhaseCountdown:
		switch k.String() {
		case "q", "esc":
			m.phase = PhaseTitle
		}

	case PhasePlaying:
		switch k.String() {
		case "up", "k", "w":
			m.keyUp = true
		case "down", "j", "s":
			m.keyDown = true
		case "p", " ":
			m.phase = PhasePaused
		case "t", "T":
			m.themeIdx = (m.themeIdx + 1) % numThemes
		case "q":
			m.phase = PhaseTitle
		case "?", "h":
			m.phase = PhaseHelp
		}
		// Release tracking (bubbletea sends key type on every event)
		switch k.Type {
		case tea.KeyUp:
			m.keyUp = true
		case tea.KeyDown:
			m.keyDown = true
		}

	case PhasePaused:
		m.keyUp = false
		m.keyDown = false
		switch k.String() {
		case "p", " ", "esc":
			m.phase = PhasePlaying
			m.lastTick = time.Now()
		case "r":
			m.startCountdown()
		case "q":
			m.phase = PhaseTitle
		}

	case PhaseGameOver:
		switch k.String() {
		case "r", "enter", " ":
			m.startCountdown()
		case "S":
			m.phase = PhaseLeaderboard
			recs, _ := m.st.LoadAll()
			m.scores = recs
		case "q", "esc":
			m.phase = PhaseTitle
		}

	case PhaseLeaderboard:
		switch k.String() {
		case "q", "esc", "enter":
			m.phase = PhaseTitle
		case "0":
			m.lbFilter = ""
			recs, _ := m.st.LoadAll()
			m.scores = recs
		case "1":
			m.lbFilter = "classic"
			m.filterScores()
		case "2":
			m.lbFilter = "arcade"
			m.filterScores()
		case "3":
			m.lbFilter = "zen"
			m.filterScores()
		case "4":
			m.lbFilter = "timed"
			m.filterScores()
		}

	case PhaseHelp:
		m.phase = PhaseTitle
	}

	return m
}

func (m *Model) filterScores() {
	all, _ := m.st.LoadAll()
	if m.lbFilter == "" {
		m.scores = all
		return
	}
	out := all[:0]
	for _, r := range all {
		if r.Mode == m.lbFilter {
			out = append(out, r)
		}
	}
	m.scores = out
}

// ─────────────────────────────────────────────────────────────────────────────
// Tick
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) tick(now time.Time) Model {
	switch m.phase {
	case PhaseCountdown:
		return m.tickCountdown(now)
	case PhasePlaying:
		return m.tickGame(now)
	}
	return m
}

func (m Model) tickCountdown(now time.Time) Model {
	if m.lastTick.IsZero() {
		m.lastTick = now
		return m
	}
	dt := now.Sub(m.lastTick).Seconds()
	m.lastTick = now

	m.cdTTL -= dt
	if m.cdTTL <= 0 {
		if m.countdown > 0 {
			m.countdown--
			m.cdTTL = 1.0
		} else {
			m.startPlaying()
		}
	}
	return m
}

func (m Model) tickGame(now time.Time) Model {
	if m.lastTick.IsZero() {
		m.lastTick = now
		return m
	}
	dt := now.Sub(m.lastTick).Seconds()
	if dt > 0.05 {
		dt = 0.05
	}
	m.lastTick = now
	m.elapsed = now.Sub(m.gameStart)

	// Time Trial
	if m.mode == ModeTimeTrial && m.elapsed >= m.timeLimit {
		m.endGame()
		return m
	}

	m.updatePaddle(dt)
	m.updateBall(dt)
	m.updateParticles(dt)
	m.updateFloatTexts(dt)

	if m.mode == ModeArcade || m.mode == ModeZen {
		m.updateFallingPUs(dt)
	}
	m.updateActivePU(dt)

	if m.bannerTTL > 0 {
		m.bannerTTL -= dt
	}
	if m.paddle.FlashTTL > 0 {
		m.paddle.FlashTTL -= dt
		if m.paddle.FlashTTL < 0 {
			m.paddle.FlashTTL = 0
		}
	}

	// Key release each frame (bubbletea doesn't give true key-up events)
	m.keyUp = false
	m.keyDown = false

	return m
}

// ─────────────────────────────────────────────────────────────────────────────
// Paddle
// ─────────────────────────────────────────────────────────────────────────────

const paddleSpeed = 28.0

func (m *Model) updatePaddle(dt float64) {
	prev := m.paddle.Y
	if m.keyUp {
		m.paddle.Y -= paddleSpeed * dt
	}
	if m.keyDown {
		m.paddle.Y += paddleSpeed * dt
	}
	m.clampPaddle()
	if dt > 0 {
		m.paddle.VY = (m.paddle.Y - prev) / dt
	}
}

func (m *Model) clampPaddle() {
	bottom := float64(m.playH-1) - float64(m.paddle.H)
	if m.paddle.Y < 1 {
		m.paddle.Y = 1
	}
	if m.paddle.Y > bottom {
		m.paddle.Y = bottom
	}
}

func (m *Model) recalcPlayArea() {
	m.playH = m.height - 6
	if m.playH < 10 {
		m.playH = 10
	}
	m.playW = m.width - 4
	if m.playW < 20 {
		m.playW = 20
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Float texts
// ─────────────────────────────────────────────────────────────────────────────

func (m *Model) updateFloatTexts(dt float64) {
	alive := m.ftexts[:0]
	for _, ft := range m.ftexts {
		ft.Life -= ft.Decay * dt
		ft.Y -= 1.5 * dt
		if ft.Life > 0 {
			alive = append(alive, ft)
		}
	}
	m.ftexts = alive
}

// ─────────────────────────────────────────────────────────────────────────────
// Active power-up timer
// ─────────────────────────────────────────────────────────────────────────────

func (m *Model) updateActivePU(dt float64) {
	if m.activePU == nil || m.activePU.Total == 0 {
		return
	}
	m.activePU.TTL -= dt
	if m.activePU.TTL <= 0 {
		m.expirePowerUp()
	}
}

func (m *Model) expirePowerUp() {
	if m.activePU == nil {
		return
	}
	switch m.activePU.Kind {
	case PUWidePaddle:
		m.paddle.H -= 3
		if m.paddle.H < m.curPhase.PaddleH {
			m.paddle.H = m.curPhase.PaddleH
		}
	case PUBomb:
		// effect was at activation, nothing to undo
	}
	m.activePU = nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Falling power-ups
// ─────────────────────────────────────────────────────────────────────────────

const puFallSpeed = 7.0

func (m *Model) updateFallingPUs(dt float64) {
	alive := m.fallingPUs[:0]
	for i := range m.fallingPUs {
		pu := m.fallingPUs[i]
		pu.Y += pu.VY * dt
		pu.FrameTTL -= dt
		if pu.FrameTTL <= 0 {
			pu.Frame = (pu.Frame + 1) % 3
			pu.FrameTTL = 0.15
		}

		px := int(pu.X + 0.5)
		py := int(pu.Y + 0.5)
		if px == PaddleX &&
			py >= int(m.paddle.Y) && py < int(m.paddle.Y)+m.paddle.H {
			m.activatePowerUp(pu.Kind)
			continue
		}
		if pu.Y >= float64(m.playH) {
			continue
		}
		alive = append(alive, pu)
	}
	m.fallingPUs = alive
}

func (m *Model) activatePowerUp(kind PowerUpKind) {
	if m.activePU != nil {
		m.expirePowerUp()
	}
	dur := kind.Duration()
	m.activePU = &ActivePU{Kind: kind, TTL: dur, Total: dur}

	switch kind {
	case PUWidePaddle:
		m.paddle.H += 3
	case PUBomb:
		m.paddle.H -= 2
		if m.paddle.H < 2 {
			m.paddle.H = 2
		}
		m.activePU.Total = 10
		m.activePU.TTL = 10
	case PUIronShield:
		m.shieldActive = true
		m.activePU.Total = 0
	case PUGhost:
		m.ghostActive = true
		m.activePU.Total = 0
	}

	m.ftexts = append(m.ftexts, FloatText{
		X:     float64(PaddleX + 3),
		Y:     m.paddle.Y,
		Text:  kind.Name(),
		Color: kind.Color(),
		Life:  1.0,
		Decay: 0.7,
	})
}

func (m *Model) spawnPowerUp() {
	kinds := []PowerUpKind{PUWidePaddle, PUSlowMo, PUFirePaddle, PUIronShield, PUGhost}
	kind := kinds[rand.Intn(len(kinds))]
	if rand.Float64() < 0.15 {
		kind = PUBomb
	}
	xMin := PaddleX + 4
	xMax := m.playW - 3
	if xMax <= xMin {
		xMax = xMin + 1
	}
	x := float64(xMin + rand.Intn(xMax-xMin))
	m.fallingPUs = append(m.fallingPUs, FallingPU{
		X:        x,
		Y:        1.0,
		Kind:     kind,
		VY:       puFallSpeed,
		FrameTTL: 0.15,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// End game
// ─────────────────────────────────────────────────────────────────────────────

func (m *Model) endGame() {
	m.phase = PhaseGameOver
	m.spawnExplosion()
	m.st.Save(store.ScoreRecord{
		Mode:        m.mode.ShortCode(),
		Score:       m.score,
		HighStreak:  m.maxStreak,
		BallsCaught: m.catches,
		BallsMissed: m.misses,
		MaxPhase:    m.curPhase.Num,
		DurationSec: int(m.elapsed.Seconds()),
		Timestamp:   time.Now(),
		Version:     "1.0.0",
	})
	if m.score > m.hiScore {
		m.hiScore = m.score
	}
}