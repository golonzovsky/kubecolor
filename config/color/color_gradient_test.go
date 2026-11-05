package color

import (
	"testing"

	"github.com/gookit/color"
)

func TestInterpolate(t *testing.T) {
	from := MustParse("#000000")
	to := MustParse("#ff00ff")
	tests := []struct {
		t    float64
		want [3]uint8
	}{
		{0, [3]uint8{0, 0, 0}},
		{0.5, [3]uint8{128, 0, 128}},
		{1, [3]uint8{255, 0, 255}},
		{-1, [3]uint8{0, 0, 0}},
		{2, [3]uint8{255, 0, 255}},
	}
	for _, tt := range tests {
		got, ok := Interpolate(from, to, tt.t)
		if !ok {
			t.Fatalf("t=%v: expected ok", tt.t)
		}
		rgb, isRGB := got.Parsed[0].(color.RGBColor)
		if !isRGB {
			t.Fatalf("t=%v: expected RGB color, got %T", tt.t, got.Parsed[0])
		}
		if [3]uint8{rgb[0], rgb[1], rgb[2]} != tt.want {
			t.Errorf("t=%v: got %v, want %v", tt.t, rgb[:3], tt.want)
		}
	}
}

func TestInterpolate_namedAndNoop(t *testing.T) {
	if _, ok := Interpolate(MustParse("red"), MustParse("blue"), 0.5); !ok {
		t.Error("named colors: expected ok")
	}
	if _, ok := Interpolate(Color{}, MustParse("blue"), 0.5); ok {
		t.Error("noop from: expected !ok")
	}
	if _, ok := Interpolate(MustParse("bold"), MustParse("blue"), 0.5); ok {
		t.Error("attribute-only from: expected !ok")
	}
	if _, ok := Interpolate(MustParse("red"), MustParse("bg=blue"), 0.5); ok {
		t.Error("background-only to: expected !ok")
	}
}
