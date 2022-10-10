package color

import (
	"github.com/muesli/termenv"
)

type Color termenv.Color

var (
	p = termenv.ColorProfile()

	// https://coolors.co/ff0054-ffc31f-73f59f-66c2cd-32a1ec-546a7b-9d4edd-7b2cbf-e0aaff-b9bfca

	Red          = p.Color("#ff0054")
	Yellow       = p.Color("#ffc31f")
	Green        = p.Color("#73f59f")
	Cyan         = p.Color("#66C2cd")
	Blue         = p.Color("#32A1EC")
	Gray         = p.Color("#546a7b")
	GrayLight    = p.Color("#b9bfca")
	Magenta      = p.Color("#9d4edd")
	MagentaDark  = p.Color("#7b2cbf")
	MagentaLight = p.Color("#e0aaff")
)

func Apply(val string, c Color) string {
	return termenv.String(val).Foreground(c).String()
}
