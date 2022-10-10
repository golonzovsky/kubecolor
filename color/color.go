package color

import (
	"github.com/muesli/termenv"
)

type Color termenv.Color

var (
	p = termenv.ColorProfile()

	Black   = p.Color("#000000")
	Gray    = p.Color("#546a7b")
	Red     = p.Color("#ff0054")
	Green   = p.Color("#73f59f")
	Yellow  = p.Color("#ffbd00")
	Blue    = p.Color("#71BEF2")
	Magenta = p.Color("#9d4edd")
	Cyan    = p.Color("#66C2CD")
	White   = p.Color("#B9BFCA")
)

func Apply(val string, c Color) string {
	return termenv.String(val).Foreground(c).String()
}
