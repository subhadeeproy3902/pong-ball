package game

import "math/rand"

// updateParticles advances every particle one timestep, removing dead ones.
func (m *Model) updateParticles(dt float64) {
	alive := m.particles[:0]
	for _, p := range m.particles {
		p.X += p.VX * dt
		p.Y += p.VY * dt
		// gravity nudge for a nicer arc
		p.VY += 3.0 * dt
		p.Life -= p.Decay * dt
		if p.Life > 0 {
			alive = append(alive, p)
		}
	}
	m.particles = alive
}

// spawnWallParticles creates the small burst when the ball hits a wall.
func (m *Model) spawnWallParticles(x, y int, side string) {
	glyphs := []rune{'·', '˙', '′', '*', '•'}
	colors := map[string]string{
		"top":    "#4ECDC4",
		"bottom": "#4ECDC4",
		"right":  "#FF6B6B",
	}
	col, ok := colors[side]
	if !ok {
		col = "#FFFFFF"
	}

	n := 4 + rand.Intn(3)
	for i := 0; i < n; i++ {
		var vx, vy float64
		switch side {
		case "right":
			vx = -(3.0 + rand.Float64()*8.0)
			vy = -4.0 + rand.Float64()*8.0
		case "top":
			vx = -4.0 + rand.Float64()*8.0
			vy = 1.0 + rand.Float64()*6.0
		case "bottom":
			vx = -4.0 + rand.Float64()*8.0
			vy = -(1.0 + rand.Float64()*6.0)
		}
		m.particles = append(m.particles, Particle{
			X:     float64(x),
			Y:     float64(y),
			VX:    vx,
			VY:    vy,
			Life:  1.0,
			Decay: 2.0 + rand.Float64()*2.0,
			Glyph: glyphs[rand.Intn(len(glyphs))],
			Color: col,
		})
	}
}