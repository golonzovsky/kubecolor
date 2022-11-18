package printer

import (
	"testing"

	"github.com/golonzovsky/kubecolor/color"
	"github.com/google/go-cmp/cmp"
)

func Test_toSpaces(t *testing.T) {
	if toSpaces(3) != "   " {
		t.Fatalf("fail")
	}
}

func Test_getColorByKeyIndent(t *testing.T) {
	tests := []struct {
		name             string
		dark             bool
		indent           int
		basicIndentWidth int
		expected         color.Color
	}{
		{"dark depth: 1", true, 2, 2, color.GrayLight},
		{"dark depth: 2", true, 4, 2, color.Yellow},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByKeyIndent(tt.indent, tt.basicIndentWidth)
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_getColorByValueType(t *testing.T) {
	tests := []struct {
		name     string
		dark     bool
		val      string
		expected color.Color
	}{
		{"dark null", true, "null", NullColor},

		{"dark bool", true, "true", BoolColor},

		{"dark number", true, "123", NumberColor},

		{"dark string", true, "aaa", StringColor},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByValueType(tt.val)
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_getColorsByBackground(t *testing.T) {
	tests := []struct {
		name     string
		dark     bool
		expected []color.Color
	}{
		{"dark", true, alternatingColors},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getAlternatingColors()
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("fail: %v", diff)
			}
		})
	}
}

func Test_getHeaderColorByBackground(t *testing.T) {
	tests := []struct {
		name     string
		dark     bool
		expected color.Color
	}{
		{"dark", true, HeaderColor},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getHeaderColor()
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_findIndent(t *testing.T) {
	tests := []struct {
		line     string
		expected int
	}{
		{"no indent", 0},
		{"  2 indent", 2},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.line, func(t *testing.T) {
			t.Parallel()
			got := findIndent(tt.line)
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}
