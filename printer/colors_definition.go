package printer

import "github.com/golonzovsky/kubecolor/color"

var (
	alternatingColors = []color.Color{
		color.BlueLight,
		color.MagentaLight,
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
