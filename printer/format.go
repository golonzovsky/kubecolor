package printer

import (
	"strconv"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// getColorByKeyIndent returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func getColorByKeyIndent(indent int, basicIndentWidth int) color.Color {
	switch indent / basicIndentWidth % 2 {
	case 1:
		return color.Blue
	default:
		return color.Blue
	}
}

// getColorByValueType returns a color by value.
// This is intended to be used to colorize any structured data e.g. Json, Yaml.
func getColorByValueType(val string) color.Color {
	if val == "null" || val == "<none>" || val == "<unknown>" {
		return NullColor
	}

	if val == "true" || val == "false" {
		return BoolColor
	}

	if _, err := strconv.Atoi(val); err == nil {
		return NumberColor
	}

	return StringColor
}

func getColorsByBackground() []color.Color {
	return colorsForDarkBackground
}

func getHeaderColor() color.Color {
	return HeaderColor
}

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
