package game

// title_art.go - PONG figlet wordmark (font: banner3) for the title screen.
// Bold block letters, ~39 cols wide so it fits any >=60-col terminal without
// wrapping; the gradient is applied at render time in view.go.

var titleWordmark = []string{
	"########   #######  ##    ##  ######   ",
	"##     ## ##     ## ###   ## ##    ##  ",
	"##     ## ##     ## ####  ## ##        ",
	"########  ##     ## ## ## ## ##   #### ",
	"##        ##     ## ##  #### ##    ##  ",
	"##        ##     ## ##   ### ##    ##  ",
	"##         #######  ##    ##  ######   ",
}
