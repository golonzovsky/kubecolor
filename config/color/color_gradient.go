package color

import (
	"github.com/gookit/color"
)

// Interpolate returns the color at fraction t between the foreground colors
// of from (t=0) and to (t=1). Returns false when either color has no
// resolvable foreground RGB value, e.g "none" or attribute-only colors.
func Interpolate(from, to Color, t float64) (Color, bool) {
	a, ok := from.fgRGB()
	if !ok {
		return Color{}, false
	}
	b, ok := to.fgRGB()
	if !ok {
		return Color{}, false
	}
	t = min(max(t, 0), 1)
	rgb := color.RGB(
		lerpUint8(a[0], b[0], t),
		lerpUint8(a[1], b[1], t),
		lerpUint8(a[2], b[2], t))
	return Color{Parsed: []ColorCode{rgb}}, true
}

func (c Color) fgRGB() (color.RGBColor, bool) {
	for _, code := range c.Parsed {
		var rgb color.RGBColor
		switch code := code.(type) {
		case color.Color:
			rgb = code.RGB()
		case color.Color256:
			rgb = code.RGB()
		case color.RGBColor:
			rgb = code
		default:
			continue
		}
		if rgb.IsEmpty() || rgb[3] == color.AsBg {
			continue
		}
		return rgb, true
	}
	return color.RGBColor{}, false
}

func lerpUint8(a, b uint8, t float64) uint8 {
	return uint8(float64(a) + (float64(b)-float64(a))*t + 0.5)
}
