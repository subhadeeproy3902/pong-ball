package game

import (
	"math"
	"math/rand"
)

// updateBall advances the ball one timestep and resolves all collisions.
func (m *Model) updateBall(dt float64) {
	// Effective speed (may be slowed by power-up)
	speedMult := 1.0
	if m.activePU != nil && m.activePU.Kind == PUSlowMo {
		speedMult = 0.65
	}

	nx := m.ball.X + m.ball.VX*speedMult*dt
	ny := m.ball.Y + m.ball.VY*speedMult*dt

	// ── Top wall ─────────────────────────────────────────────────────────
	if ny <= 1 {
		ny = 2 - ny
		m.ball.VY = -m.ball.VY
		m.spawnWallParticles(int(nx), 1, "top")
	}

	// ── Bottom wall ───────────────────────────────────────────────────────
	if ny >= float64(m.playH-2) {
		ny = float64(m.playH-2)*2 - ny
		m.ball.VY = -m.ball.VY
		m.spawnWallParticles(int(nx), m.playH-2, "bottom")
	}

	// ── Right wall ────────────────────────────────────────────────────────
	rightEdge := float64(m.playW - 2)
	if nx >= rightEdge {
		nx = rightEdge*2 - nx
		m.ball.VX = -m.ball.VX
		m.spawnWallParticles(m.playW-2, int(ny), "right")
	}

	// ── Paddle collision ──────────────────────────────────────────────────
	paddleRight := float64(PaddleX + 1)

	if m.ball.VX < 0 && nx <= paddleRight && m.ball.X > paddleRight {
		paddleTop := m.paddle.Y
		paddleBot := m.paddle.Y + float64(m.paddle.H)

		if ny >= paddleTop && ny <= paddleBot {
			m.handlePaddleHit(&nx, &ny, ny, paddleTop, paddleBot)
		} else {
			// Missed the paddle
			m.handleMiss()
			return
		}
	}

	// ── Trail update ──────────────────────────────────────────────────────
	trailLen := m.curPhase.TrailLen
	m.ball.Trail = append([]Pt{{X: int(math.Round(m.ball.X)), Y: int(math.Round(m.ball.Y))}}, m.ball.Trail...)
	if len(m.ball.Trail) > trailLen {
		m.ball.Trail = m.ball.Trail[:trailLen]
	}

	m.ball.X = nx
	m.ball.Y = ny
}

// handlePaddleHit calculates the new velocity vector after a paddle hit.
func (m *Model) handlePaddleHit(nx, ny *float64, ballY, paddleTop, paddleBot float64) {
	// Ghost power-up: pass through once
	if m.ghostActive {
		m.ghostActive = false
		if m.activePU != nil && m.activePU.Kind == PUGhost {
			m.activePU = nil
		}
		return
	}

	paddleH := paddleBot - paddleTop
	relHit := ((ballY - paddleTop) / paddleH) - 0.5 // −0.5 to +0.5
	maxAngle := math.Pi / 3                          // 60°
	angle := relHit * maxAngle

	speed := math.Sqrt(m.ball.VX*m.ball.VX+m.ball.VY*m.ball.VY)

	// Phase speed escalation
	speed = math.Max(speed, BaseSpeed*m.curPhase.SpeedMult)

	m.ball.VX = speed * math.Cos(angle)
	m.ball.VY = speed * math.Sin(angle)

	// Spin transfer from paddle movement
	m.ball.VY += m.paddle.VY * 0.35

	// Ensure ball is moving right (away from paddle)
	if m.ball.VX <= 0 {
		m.ball.VX = math.Abs(m.ball.VX) + 1
	}

	*nx = float64(PaddleX+1) + 0.5

	// Edge bonus detection
	edgeFraction := 0.12
	isEdge := relHit > (0.5-edgeFraction) || relHit < (-0.5+edgeFraction)

	m.scoreHit(isEdge)

	// Visual flash
	m.paddle.FlashTTL = 0.12

	// Phase transition check
	m.checkPhaseTransition()

	// Power-up spawn (Arcade / Zen every 7 catches)
	if m.mode == ModeArcade || m.mode == ModeZen {
		m.catchesSinceLastPU++
		if m.catchesSinceLastPU >= 7 {
			m.catchesSinceLastPU = 0
			m.spawnPowerUp()
		}
	}
}

// handleMiss is called when the ball passes the paddle.
func (m *Model) handleMiss() {
	m.misses++
	m.streak = 0

	// Iron shield: one auto-save
	if m.shieldActive {
		m.shieldActive = false
		if m.activePU != nil && m.activePU.Kind == PUIronShield {
			m.activePU = nil
		}
		// Bounce ball back
		m.ball.VX = math.Abs(m.ball.VX)
		m.ball.X = float64(PaddleX + 2)
		return
	}

	if m.mode == ModeZen {
		// No lives lost; just reset ball position
		m.resetBallOnly()
		return
	}

	m.lives--
	if m.lives <= 0 {
		m.endGame()
	} else {
		m.resetBallOnly()
	}
}

func (m *Model) resetBallOnly() {
	ph := m.curPhase
	speed := BaseSpeed * ph.SpeedMult
	angle := -0.3 + rand.Float64()*0.6
	m.ball = Ball{
		X:  float64(m.playW / 2),
		Y:  float64(m.playH / 2),
		VX: math.Abs(speed * math.Cos(angle)),
		VY: speed * math.Sin(angle),
	}
}

// spawnExplosion creates the game-over particle burst at the ball's position.
func (m *Model) spawnExplosion() {
	glyphs := []rune{'✦', '✧', '★', '✶', '✷', '✸', '✹', '✺', '·', '*'}
	colors := []string{"#FFD700", "#FF5370", "#00FFFF", "#C3E88D", "#FFCB6B", "#89DDFF"}
	for i := 0; i < 16; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 4.0 + rand.Float64()*10.0
		m.particles = append(m.particles, Particle{
			X:     m.ball.X,
			Y:     m.ball.Y,
			VX:    speed * math.Cos(angle),
			VY:    speed * math.Sin(angle),
			Life:  1.0,
			Decay: 0.5 + rand.Float64()*0.5,
			Glyph: glyphs[rand.Intn(len(glyphs))],
			Color: colors[rand.Intn(len(colors))],
		})
	}
}