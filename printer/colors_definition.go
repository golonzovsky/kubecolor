package printer

import "github.com/kubecolor/kubecolor/color"

var (
	// Preset of colors for background
	// Please use them when you just need random colors
	colorsForDarkBackground = []color.Color{
		color.GrayLight,
		color.Cyan,
	}

	// colors to be recommended to be used for some context
	// e.g. Json, Yaml, kubectl-describe format etc.

	KeyColor    = color.Cyan
	StringColor = color.Green
	BoolColor   = color.MagentaLight
	NumberColor = color.Magenta
	NullColor   = color.Gray
	HeaderColor = color.GrayLight // for plain table
)
