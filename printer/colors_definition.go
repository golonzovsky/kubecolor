package printer

import "github.com/kubecolor/kubecolor/color"

var (
	// Preset of colors for background
	// Please use them when you just need random colors
	colorsForDarkBackground = []color.Color{
		color.White,
		color.Cyan,
	}

	colorsForLightBackground = []color.Color{
		color.Black,
		color.Blue,
	}

	// colors to be recommended to be used for some context
	// e.g. Json, Yaml, kubectl-describe format etc.

	KeyColorForDark    = color.Cyan
	StringColorForDark = color.White
	BoolColorForDark   = color.Green
	NumberColorForDark = color.Magenta
	NullColorForDark   = color.Yellow
	HeaderColorForDark = color.White // for plain table
)
