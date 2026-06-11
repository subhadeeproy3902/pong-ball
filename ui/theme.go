package ui

import "github.com/charmbracelet/lipgloss"

// ThemeCount is exported so the game package can cycle themes.
const ThemeCount = 4

// Theme holds the full color palette for one visual style.
type Theme struct {
	Name        string
	Paddle      string
	Ball        string
	Trail       []string // index 0 = farthest from ball, last = nearest
	TopBot      string   // top/bottom wall color
	RightWall   string   // right wall color
	Score       string
	HiScore     string
	Streak      string
	Banner      string
	Background  string
	Border      string
	HeaderBg    string
	PowerUp     string
	LivesColor  string
}

// Themes contains all four built-in themes.
var Themes = []Theme{
	{
		// ── 0 · Neon Arcade (default) ─────────────────────────────────────
		Name:       "Neon",
		Paddle:     "#00FFFF",
		Ball:       "#FFD700",
		Trail:      []string{"#8B0000", "#FF4500", "#FF8C00"},
		TopBot:     "#4ECDC4",
		RightWall:  "#FF6B6B",
		Score:      "#C3E88D",
		HiScore:    "#FFCB6B",
		Streak:     "#FF5370",
		Banner:     "#FF00FF",
		Background: "#0D0D1A",
		Border:     "#2D2D44",
		HeaderBg:   "#12122A",
		PowerUp:    "#89DDFF",
		LivesColor: "#FF5370",
	},
	{
		// ── 1 · Monochrome ────────────────────────────────────────────────
		Name:       "Mono",
		Paddle:     "#FFFFFF",
		Ball:       "#FFFFFF",
		Trail:      []string{"#333333", "#666666", "#999999"},
		TopBot:     "#AAAAAA",
		RightWall:  "#AAAAAA",
		Score:      "#CCCCCC",
		HiScore:    "#FFFFFF",
		Streak:     "#FFFFFF",
		Banner:     "#FFFFFF",
		Background: "#000000",
		Border:     "#333333",
		HeaderBg:   "#0A0A0A",
		PowerUp:    "#CCCCCC",
		LivesColor: "#FFFFFF",
	},
	{
		// ── 2 · Sunset ────────────────────────────────────────────────────
		Name:       "Sunset",
		Paddle:     "#FFA07A",
		Ball:       "#FFD700",
		Trail:      []string{"#4A0000", "#8B1A1A", "#CD5C5C"},
		TopBot:     "#FF6347",
		RightWall:  "#FF4500",
		Score:      "#FFD700",
		HiScore:    "#FFF8DC",
		Streak:     "#FF6B6B",
		Banner:     "#FF8C00",
		Background: "#1A0800",
		Border:     "#4A2000",
		HeaderBg:   "#200A00",
		PowerUp:    "#FFD700",
		LivesColor: "#FF6347",
	},
	{
		// ── 3 · Ocean Night ───────────────────────────────────────────────
		Name:       "Ocean",
		Paddle:     "#90E0EF",
		Ball:       "#CAF0F8",
		Trail:      []string{"#03045E", "#0077B6", "#00B4D8"},
		TopBot:     "#0096C7",
		RightWall:  "#0077B6",
		Score:      "#90E0EF",
		HiScore:    "#CAF0F8",
		Streak:     "#48CAE4",
		Banner:     "#ADE8F4",
		Background: "#000814",
		Border:     "#023E8A",
		HeaderBg:   "#001D3D",
		PowerUp:    "#ADE8F4",
		LivesColor: "#48CAE4",
	},
}

// ThemeIndexByName returns the index of the named theme, defaulting to 0.
func ThemeIndexByName(name string) int {
	for i, t := range Themes {
		if t.Name == name {
			return i
		}
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Style factories (called with the current theme)
// ─────────────────────────────────────────────────────────────────────────────

func Color(hex string) lipgloss.Color { return lipgloss.Color(hex) }

func Styled(hex string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(Color(hex))
}

func Bold(hex string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(Color(hex)).Bold(true)
}

func BgFg(bg, fg string) lipgloss.Style {
	return lipgloss.NewStyle().Background(Color(bg)).Foreground(Color(fg))
}

// Box creates a lipgloss border box in theme border color.
func Box(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Color(t.Border)).
		Padding(0, 1)
}