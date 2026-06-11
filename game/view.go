package game

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/subhadeeproy3902/paddle-ball/ui"
)

// View is the bubbletea view function — called every frame.
func (m Model) View() string {
	if m.width < 80 || m.height < 24 {
		return m.viewTooSmall()
	}

	t := ui.Themes[m.themeIdx]

	switch m.phase {
	case PhaseTitle:
		return m.viewTitle(t)
	case PhaseCountdown:
		return m.viewCountdown(t)
	case PhasePlaying:
		return m.viewPlaying(t)
	case PhasePaused:
		return m.viewPaused(t)
	case PhaseGameOver:
		return m.viewGameOver(t)
	case PhaseLeaderboard:
		return m.viewLeaderboard(t)
	case PhaseHelp:
		return m.viewHelp(t)
	}
	return ""
}

// ─────────────────────────────────────────────────────────────────────────────
// Too-small guard
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewTooSmall() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5370")).
		Bold(true).
		Render(fmt.Sprintf(
			"\n  Terminal too small!\n  Minimum: 80×24\n  Current: %d×%d\n  Please resize.",
			m.width, m.height,
		))
}

// ─────────────────────────────────────────────────────────────────────────────
// Title screen
// ─────────────────────────────────────────────────────────────────────────────

const logo = `
 ██████╗  █████╗ ██████╗ ██████╗ ██╗     ███████╗
 ██╔══██╗██╔══██╗██╔══██╗██╔══██╗██║     ██╔════╝
 ██████╔╝███████║██║  ██║██║  ██║██║     █████╗  
 ██╔═══╝ ██╔══██║██║  ██║██║  ██║██║     ██╔══╝  
 ██║     ██║  ██║██████╔╝██████╔╝███████╗███████╗
 ╚═╝     ╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚══════╝╚══════╝
  ██████╗  █████╗ ██╗     ██╗                     
  ██╔══██╗██╔══██╗██║     ██║                     
  ██████╔╝███████║██║     ██║                     
  ██╔══██╗██╔══██║██║     ██║                     
  ██████╔╝██║  ██║███████╗███████╗                
  ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝`

func (m Model) viewTitle(t ui.Theme) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Paddle))
	logoStr := logoStyle.Render(logo)

	modes := []string{"  Classic   · One life, pure score chase",
		"  Arcade    · 3 lives + power-ups",
		"  Zen       · Infinite lives, just vibe",
		"  Time Trial · 60-second blitz"}

	var menu strings.Builder
	for i, line := range modes {
		num := fmt.Sprintf("[%d] ", i+1)
		if i == m.menuSel {
			s := lipgloss.NewStyle().
				Foreground(lipgloss.Color(t.Paddle)).
				Bold(true).
				Background(lipgloss.Color("#1A1A2E"))
			menu.WriteString(s.Render("▶ "+num+line[4:]) + "\n")
		} else {
			s := lipgloss.NewStyle().Foreground(lipgloss.Color("#666688"))
			menu.WriteString(s.Render("  "+num+line[4:]) + "\n")
		}
	}

	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Border)).
		Padding(0, 2).
		Width(46).
		Render(menu.String() +
			"\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#555577")).Render("  [S] Scores  [T] Theme  [?] Help  [Q] Quit"))

	hiStr := lipgloss.NewStyle().Foreground(lipgloss.Color(t.HiScore)).
		Render(fmt.Sprintf("Hi-Score: %d", m.hiScore))
	themeStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#555577")).
		Render(fmt.Sprintf("Theme: %s", t.Name))

	footer := lipgloss.JoinHorizontal(lipgloss.Top,
		hiStr,
		strings.Repeat(" ", 4),
		themeStr,
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, logoStr, "\n", menuBox, "\n", footer),
	)
}

// ─────────────────────────────────────────────────────────────────────────────
// Countdown screen
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewCountdown(t ui.Theme) string {
	var label string
	if m.countdown == 0 {
		label = "GO!"
	} else {
		label = fmt.Sprintf("%d", m.countdown)
	}
	big := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Paddle)).
		Bold(true).
		Width(10).
		Align(lipgloss.Center).
		Render(label)

	mode := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888899")).
		Render("Mode: " + m.mode.String())

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, big, "\n", mode),
	)
}

// ─────────────────────────────────────────────────────────────────────────────
// Playing screen
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewPlaying(t ui.Theme) string {
	header := m.buildHeader(t)
	playArea := m.buildPlayArea(t)
	footer := m.buildFooter(t)

	// Phase banner overlay
	var bannerRow string
	if m.bannerTTL > 0 {
		bs := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.bannerColor)).
			Bold(true).
			Width(m.width).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#0A0A1A"))
		bannerRow = bs.Render(m.bannerText) + "\n"
	}

	return header + "\n" + bannerRow + playArea + "\n" + footer
}

func (m Model) buildHeader(t ui.Theme) string {
	title := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Paddle)).Bold(true).Render("🏓 PADDLEBALL")
	modeS := lipgloss.NewStyle().Foreground(lipgloss.Color("#666688")).Render("Mode: " + m.mode.String())

	scoreS := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Score)).Bold(true).
		Render(fmt.Sprintf("SCORE: %03d", m.score))
	hiS := lipgloss.NewStyle().Foreground(lipgloss.Color(t.HiScore)).
		Render(fmt.Sprintf("🏆 HI: %d", m.hiScore))

	left := title + "  " + modeS
	right := scoreS + "  " + hiS

	gap := m.width - visLen(left) - visLen(right) - 2
	if gap < 1 {
		gap = 1
	}

	row1 := left + strings.Repeat(" ", gap) + right

	// Row 2: lives / streak / phase
	var livesStr string
	if m.mode == ModeArcade {
		hearts := strings.Repeat("❤️  ", m.lives)
		livesStr = lipgloss.NewStyle().Foreground(lipgloss.Color(t.LivesColor)).Render(hearts)
	}

	streakStr := ""
	if m.streak >= 5 {
		streakStr = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Streak)).Bold(true).
			Render(fmt.Sprintf("🔥 ×%d", m.streak))
	}

	phaseStr := lipgloss.NewStyle().Foreground(lipgloss.Color(m.curPhase.Color)).
		Render(m.curPhase.Emoji + " " + m.curPhase.Name)

	// Time Trial countdown
	var timerStr string
	if m.mode == ModeTimeTrial {
		remain := m.timeLimit - m.elapsed
		if remain < 0 {
			remain = 0
		}
		rem := remain.Round(time.Second)
		secs := int(rem.Seconds())
		tc := t.Score
		if secs <= 10 {
			tc = t.Streak
		}
		timerStr = lipgloss.NewStyle().Foreground(lipgloss.Color(tc)).Bold(true).
			Render(fmt.Sprintf("⏱  %ds", secs))
	}

	row2Parts := []string{}
	if livesStr != "" {
		row2Parts = append(row2Parts, livesStr)
	}
	if streakStr != "" {
		row2Parts = append(row2Parts, streakStr)
	}
	if timerStr != "" {
		row2Parts = append(row2Parts, timerStr)
	}
	row2Parts = append(row2Parts, phaseStr)
	row2 := strings.Join(row2Parts, "   ")

	style := lipgloss.NewStyle().
		Background(lipgloss.Color(t.HeaderBg)).
		Width(m.width).
		Padding(0, 1)
	return style.Render(row1) + "\n" + style.Render(row2)
}

func (m Model) buildFooter(t ui.Theme) string {
	// Power-up progress bar
	var puStr string
	if m.activePU != nil && m.activePU.Total > 0 {
		pct := m.activePU.TTL / m.activePU.Total
		barW := 10
		filled := int(float64(barW) * pct)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barW-filled)
		puStr = lipgloss.NewStyle().Foreground(lipgloss.Color(m.activePU.Kind.Color())).
			Render(fmt.Sprintf("%s %s %.0fs", string(m.activePU.Kind.Glyph()), bar, m.activePU.TTL))
		puStr += "   "
	} else if m.activePU != nil {
		puStr = lipgloss.NewStyle().Foreground(lipgloss.Color(m.activePU.Kind.Color())).
			Render(string(m.activePU.Kind.Glyph()) + " " + m.activePU.Kind.Name() + " (active)") + "   "
	}

	controls := lipgloss.NewStyle().Foreground(lipgloss.Color("#444466")).
		Render("[↑↓·WS·KJ] Move  [P] Pause  [T] Theme  [?] Help  [Q] Quit")

	row := puStr + controls

	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.HeaderBg)).
		Width(m.width).
		Padding(0, 1).
		Render(row)
}

// ─────────────────────────────────────────────────────────────────────────────
// Play area renderer
// ─────────────────────────────────────────────────────────────────────────────

type cell struct {
	r     rune
	color string
}

func (m Model) buildPlayArea(t ui.Theme) string {
	grid := make([][]cell, m.playH)
	for y := range grid {
		grid[y] = make([]cell, m.playW)
		for x := range grid[y] {
			grid[y][x] = cell{' ', ""}
		}
	}

	// Top wall
	for x := 0; x < m.playW-1; x++ {
		grid[0][x] = cell{'─', t.TopBot}
	}
	// Bottom wall
	for x := 0; x < m.playW-1; x++ {
		grid[m.playH-1][x] = cell{'─', t.TopBot}
	}
	// Right wall
	for y := 0; y < m.playH; y++ {
		grid[y][m.playW-1] = cell{'┃', t.RightWall}
	}
	grid[0][m.playW-1] = cell{'┐', t.RightWall}
	grid[m.playH-1][m.playW-1] = cell{'┘', t.RightWall}

	// Ball trail
	trailRunes := []rune{'░', '▒', '▓'}
	for i, pt := range m.ball.Trail {
		if pt.Y >= 0 && pt.Y < m.playH && pt.X >= 0 && pt.X < m.playW {
			ri := len(m.ball.Trail) - 1 - i
			if ri >= 0 && ri < len(trailRunes) {
				col := t.Trail[len(t.Trail)-1]
				if ri < len(t.Trail) {
					col = t.Trail[ri]
				}
				grid[pt.Y][pt.X] = cell{trailRunes[ri], col}
			}
		}
	}

	// Ball
	bx := int(math.Round(m.ball.X))
	by := int(math.Round(m.ball.Y))
	if m.phase != PhaseGameOver && by >= 0 && by < m.playH && bx >= 0 && bx < m.playW {
		grid[by][bx] = cell{'●', t.Ball}
	}

	// Paddle
	py := int(m.paddle.Y)
	padColor := t.Paddle
	if m.paddle.FlashTTL > 0 {
		padColor = "#FFFFFF"
	}
	if m.activePU != nil && m.activePU.Kind == PUFirePaddle {
		padColor = "#FF8C00"
	}
	if m.shieldActive {
		padColor = "#4ECDC4"
	}
	for dy := 0; dy < m.paddle.H; dy++ {
		y := py + dy
		if y >= 1 && y < m.playH-1 && PaddleX < m.playW {
			grid[y][PaddleX] = cell{'█', padColor}
		}
	}

	// Shield glyph (left-wall protection)
	if m.shieldActive && PaddleX-1 >= 0 {
		for dy := 0; dy < m.paddle.H; dy++ {
			y := py + dy
			if y >= 1 && y < m.playH-1 {
				grid[y][0] = cell{'╟', "#4ECDC4"}
			}
		}
	}

	// Particles
	for _, p := range m.particles {
		px := int(math.Round(p.X))
		ppy := int(math.Round(p.Y))
		if ppy >= 0 && ppy < m.playH && px >= 0 && px < m.playW {
			grid[ppy][px] = cell{p.Glyph, p.Color}
		}
	}

	// Falling power-ups
	frames := []rune{'▿', '▾', '▽'}
	for _, pu := range m.fallingPUs {
		fpx := int(math.Round(pu.X))
		fpy := int(math.Round(pu.Y))
		if fpy >= 0 && fpy < m.playH && fpx >= 0 && fpx < m.playW {
			fr := frames[pu.Frame%3]
			col := pu.Kind.Color()
			if fpy%2 == 0 {
				fr = pu.Kind.Glyph()
			}
			grid[fpy][fpx] = cell{fr, col}
		}
	}

	// Float texts
	for _, ft := range m.ftexts {
		fx := int(ft.X)
		fy := int(ft.Y)
		for i, ch := range ft.Text {
			tx := fx + i
			if fy >= 0 && fy < m.playH && tx >= 0 && tx < m.playW {
				grid[fy][tx] = cell{ch, ft.Color}
			}
		}
	}

	// Render to string
	var sb strings.Builder
	for _, row := range grid {
		for _, c := range row {
			if c.color != "" {
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.color)).Render(string(c.r)))
			} else {
				sb.WriteRune(c.r)
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// ─────────────────────────────────────────────────────────────────────────────
// Pause overlay
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewPaused(t ui.Theme) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(t.Paddle)).
		Padding(1, 4)

	dur := m.elapsed.Round(time.Second)
	content := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Paddle)).Bold(true).Render("⏸  PAUSED") +
		"\n\n" +
		fmt.Sprintf("Score so far: %d\n", m.score) +
		fmt.Sprintf("Streak:       ×%d\n", m.streak) +
		fmt.Sprintf("Time played:  %s\n", fmtElapsed(dur)) +
		"\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888899")).Render(
			"[P / Space]  Resume\n[R]          Restart\n[Q]          Quit to title")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
}

// ─────────────────────────────────────────────────────────────────────────────
// Game over
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewGameOver(t ui.Theme) string {
	rank, rankColor := RankForScore(m.score)
	newBest := m.score > 0 && m.score == m.hiScore

	titleLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5370")).Bold(true).
		Width(44).Align(lipgloss.Center).Render("GAME  OVER")

	row := func(label, val string) string {
		return fmt.Sprintf("  %-18s│  %s\n", label, val)
	}

	scoreStr := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Score)).Bold(true).Render(fmt.Sprintf("%d", m.score))
	rankStr := lipgloss.NewStyle().Foreground(lipgloss.Color(rankColor)).Bold(true).Render(rank)
	streakStr := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Streak)).Render(fmt.Sprintf("×%d", m.maxStreak))

	stats := row("Final Score", scoreStr) +
		row("Rank", rankStr) +
		row("High Streak", streakStr) +
		row("Balls Caught", fmt.Sprintf("%d", m.catches)) +
		row("Balls Missed", fmt.Sprintf("%d", m.misses)) +
		row("Max Phase", fmt.Sprintf("%s %s", m.curPhase.Emoji, m.curPhase.Name)) +
		row("Time Played", fmtElapsed(m.elapsed.Round(time.Second)))

	var bestLine string
	if newBest {
		bestLine = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true).
			Width(44).Align(lipgloss.Center).Render("★  NEW PERSONAL BEST!  ★")
	}

	actions := lipgloss.NewStyle().Foreground(lipgloss.Color("#666688")).
		Render("\n  [R / Enter]  Play Again   [S] Scores   [Q] Quit")

	boxContent := titleLine + "\n" +
		strings.Repeat("─", 44) + "\n" +
		stats +
		strings.Repeat("─", 44) +
		bestLine +
		actions

	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(t.Border)).
		Padding(0, 1).
		Render(boxContent)

	// Play area behind (with particles)
	bg := m.buildPlayArea(t)
	_ = bg // could overlay, but lipgloss.Place is simpler

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// ─────────────────────────────────────────────────────────────────────────────
// Leaderboard
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewLeaderboard(t ui.Theme) string {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.HiScore)).Bold(true)
	title := titleStyle.Render("🏆  PADDLEBALL — SCORE HISTORY")

	filterLabel := "All Modes"
	if m.lbFilter != "" {
		filterLabel = strings.Title(m.lbFilter)
	}
	sub := lipgloss.NewStyle().Foreground(lipgloss.Color("#666688")).Render("Showing: " + filterLabel)

	// Table header
	hdr := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Paddle)).Bold(true).
		Render("  #   SCORE   MODE       STREAK   DURATION   DATE")
	sep := strings.Repeat("─", 54)

	rows := []string{hdr, sep}
	display := m.scores
	if len(display) > 12 {
		display = display[:12]
	}
	for i, r := range display {
		rank, _ := RankForScore(r.Score)
		durS := fmtSecs(r.DurationSec)
		line := fmt.Sprintf("  %-3d %-7d %-10s ×%-7d %-10s %s",
			i+1, r.Score, r.Mode, r.HighStreak, durS, r.Timestamp.Format("Jan 02"))
		var col string
		if i == 0 {
			col = "#FFD700"
		} else if i == 1 {
			col = "#AAAACC"
		} else if i == 2 {
			col = "#CD7F32"
		} else {
			col = "#666688"
		}
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(col)).Render(line))
		_ = rank
	}

	if len(m.scores) == 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color("#555577")).Render("  No scores yet. Play a game!"))
	}

	rows = append(rows, sep)
	stats := m.st.Aggregate(m.scores)
	statsLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#555577")).
		Render(fmt.Sprintf("  Caught: %d total · Played: %s · Best streak: ×%d",
			stats.TotalCaught, fmtSecs(stats.TotalTimeSec), stats.BestStreak))
	rows = append(rows, statsLine)
	rows = append(rows, "")
	rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color("#444466")).
		Render("  [0] All  [1] Classic  [2] Arcade  [3] Zen  [4] Timed   [Q] Back"))

	table := strings.Join(rows, "\n")
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Border)).
		Padding(1, 2).
		Render(table)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, sub, "\n", box))
}

// ─────────────────────────────────────────────────────────────────────────────
// Help screen
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewHelp(t ui.Theme) string {
	header := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Paddle)).Bold(true).Render("🏓  PADDLEBALL — CONTROLS")

	col := func(k, v string) string {
		return fmt.Sprintf("  %-22s%s\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color(t.Score)).Render(k),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#888899")).Render(v))
	}

	content :=
		col("↑ / W / K", "Move paddle up") +
			col("↓ / S / J", "Move paddle down") +
			col("P / Space", "Pause / Resume") +
			col("T", "Cycle color theme") +
			col("?", "Toggle this help") +
			col("R", "Restart (from pause/game over)") +
			col("Q / Ctrl+C", "Quit") +
			col("1–4 (title)", "Select game mode") +
			"\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00")).Bold(true).Render("  POWER-UPS (Arcade / Zen)") + "\n" +
			col("Ⓦ  Wide Paddle", "Paddle grows +3 cells for 12s") +
			col("ⓢ  Slow Mo", "Ball slows 35% for 8s") +
			col("ⓕ  Fire Paddle", "Score ×2 per hit for 15s") +
			col("ⓘ  Iron Shield", "One auto-save bounce (one-time)") +
			col("ⓖ  Ghost Ball", "Pass-through once, no penalty") +
			col("Ⓑ  BOMB (dodge!)", "Paddle shrinks −2 cells for 10s") +
			"\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#555577")).Render("  Press any key to go back")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Border)).
		Padding(1, 3).
		Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, header, "\n", box))
}

// ─────────────────────────────────────────────────────────────────────────────
// String helpers
// ─────────────────────────────────────────────────────────────────────────────

func fmtElapsed(d time.Duration) string {
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	if m > 0 {
		return fmt.Sprintf("%dm %02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func fmtSecs(sec int) string {
	d := time.Duration(sec) * time.Second
	return fmtElapsed(d)
}

// visLen approximates the visible (non-ANSI) length of a string.
func visLen(s string) int {
	inEsc := false
	count := 0
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		count += utf8.RuneLen(r)
	}
	return count
}