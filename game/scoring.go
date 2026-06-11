package game

import (
	"fmt"
	"math"
)

// scoreHit awards points for a successful paddle hit.
func (m *Model) scoreHit(edgeHit bool) {
	m.catches++
	m.streak++
	if m.streak > m.maxStreak {
		m.maxStreak = m.streak
	}

	// Base points
	pts := 1
	if edgeHit {
		pts += 2
	}

	// Phase bonus
	switch m.curPhase.Num {
	case 4:
		pts += 2
	case 5:
		pts += 4
	}

	// Streak multiplier
	mult := m.streakMultiplier()

	// Fire paddle power-up multiplier
	if m.activePU != nil && m.activePU.Kind == PUFirePaddle {
		mult *= 2.0
	}

	earned := int(float64(pts) * mult)
	m.score += earned

	// Rally bonus milestones
	bonus := 0
	if m.streak == 50 {
		bonus = 25
	} else if m.streak == 100 {
		bonus = 75
	}
	if bonus > 0 {
		m.score += bonus
		earned += bonus
	}

	// Float text
	label := fmt.Sprintf("+%d", earned)
	color := scoreColor(earned)
	if mult > 1.0 {
		label += fmt.Sprintf(" ×%.1g", mult)
	}
	m.ftexts = append(m.ftexts, FloatText{
		X:     float64(PaddleX + 3),
		Y:     m.paddle.Y + float64(m.paddle.H/2),
		Text:  label,
		Color: color,
		Life:  1.0,
		Decay: 1.2,
	})
}

// streakMultiplier returns the current score multiplier based on streak length.
func (m *Model) streakMultiplier() float64 {
	switch {
	case m.streak >= 35:
		return 3.0
	case m.streak >= 20:
		return 2.0
	case m.streak >= 10:
		return 1.5
	default:
		return 1.0
	}
}

// checkPhaseTransition upgrades difficulty when the score crosses a threshold.
func (m *Model) checkPhaseTransition() {
	newPhase := PhaseForScore(m.score)
	if newPhase.Num > m.curPhase.Num {
		old := m.curPhase
		m.curPhase = newPhase
		m.paddle.H = newPhase.PaddleH
		_ = old

		m.bannerText = fmt.Sprintf("  %s  PHASE %d — %s  %s  ",
			newPhase.Emoji, newPhase.Num, newPhase.Name, newPhase.Emoji)
		m.bannerColor = newPhase.Color
		m.bannerTTL = 1.8

		// Bump ball speed
		speed := BaseSpeed * newPhase.SpeedMult
		curr := magnitude(m.ball.VX, m.ball.VY)
		if curr > 0 {
			m.ball.VX = m.ball.VX / curr * speed
			m.ball.VY = m.ball.VY / curr * speed
		}
	}
}

// magnitude returns the Euclidean magnitude of a 2D vector.
func magnitude(vx, vy float64) float64 {
	return math.Sqrt(vx*vx + vy*vy)
}

// RankForScore returns a rank title for a final score.
func RankForScore(score int) (string, string) {
	switch {
	case score >= 500:
		return "🏆 God Mode", "#FFD700"
	case score >= 200:
		return "💎 Grandmaster", "#89DDFF"
	case score >= 100:
		return "🌟 Legend", "#C3E88D"
	case score >= 50:
		return "⚡ Speedster", "#FFCB6B"
	case score >= 25:
		return "🔥 Baller", "#FF8C00"
	case score >= 10:
		return "🎮 Rookie", "#FF5370"
	default:
		return "🐣 Hatchling", "#AAAAAA"
	}
}

func scoreColor(pts int) string {
	switch {
	case pts >= 10:
		return "#FF5370"
	case pts >= 5:
		return "#FFCB6B"
	case pts >= 3:
		return "#C3E88D"
	default:
		return "#00FFFF"
	}
}