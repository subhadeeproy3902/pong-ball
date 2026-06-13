package game

// title_art.go - big "Paddle Ball" wordmark for the title screen. Shown when
// the terminal is wide enough; view.go falls back to a compact wordmark on
// narrow terminals so it never wraps. Gradient is applied at render time.

var titleWordmark = []string{
	" /$$$$$$$                 /$$       /$$ /$$                 /$$$$$$$            /$$ /$$",
	"| $$__  $$               | $$      | $$| $$                | $$__  $$          | $$| $$",
	"| $$  \\ $$ /$$$$$$   /$$$$$$$  /$$$$$$$| $$  /$$$$$$       | $$  \\ $$  /$$$$$$ | $$| $$",
	"| $$$$$$$/|____  $$ /$$__  $$ /$$__  $$| $$ /$$__  $$      | $$$$$$$  |____  $$| $$| $$",
	"| $$____/  /$$$$$$$| $$  | $$| $$  | $$| $$| $$$$$$$$      | $$__  $$  /$$$$$$$| $$| $$",
	"| $$      /$$__  $$| $$  | $$| $$  | $$| $$| $$_____/      | $$  \\ $$ /$$__  $$| $$| $$",
	"| $$     |  $$$$$$$|  $$$$$$$|  $$$$$$$| $$|  $$$$$$$      | $$$$$$$/|  $$$$$$$| $$| $$",
	"|__/      \\_______/ \\_______/ \\_______/|__/ \\_______/      |_______/  \\_______/|__/|__/",
}
