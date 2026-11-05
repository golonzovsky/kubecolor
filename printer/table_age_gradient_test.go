package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
)

func Test_AgeGradientTablePrinter(t *testing.T) {
	// blue -> red: the midpoint is magenta-ish, distinct from both ends at any color level
	theme := &config.Theme{
		Data: config.ThemeData{
			DurationFresh: color.MustParseSlice("#0000ff"),
			Duration:      color.MustParse("#ff0000"),
		},
	}
	input := "NAME  AGE\na     1m\nb     1h\nc     1d\n"
	p := &AgeGradientTablePrinter{Table: NewTablePrinter(true, theme, nil), Theme: theme}
	var buf bytes.Buffer
	p.Print(strings.NewReader(input), &buf)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %q", len(lines), buf.String())
	}
	newest := color.MustParse("#0000ff").Render("1m")
	oldest := color.MustParse("#ff0000").Render("1d")
	if !strings.HasSuffix(lines[1], newest) {
		t.Errorf("newest: expected suffix %q, got %q", newest, lines[1])
	}
	if !strings.HasSuffix(lines[3], oldest) {
		t.Errorf("oldest: expected suffix %q, got %q", oldest, lines[3])
	}
	mid := lines[2]
	if !strings.Contains(mid, "\x1b[") || !strings.HasSuffix(mid, "m1h\x1b[0m") {
		t.Errorf("middle: expected colored 1h, got %q", mid)
	}
	if strings.Contains(mid, color.MustParse("#0000ff").ANSICode()) || strings.Contains(mid, color.MustParse("#ff0000").ANSICode()) {
		t.Errorf("middle: expected a color between the ends, got %q", mid)
	}
}

func Test_AgeGradientTablePrinter_fallsBackWithoutRange(t *testing.T) {
	theme := &config.Theme{
		Data: config.ThemeData{
			DurationFresh: color.MustParseSlice("green"),
			Duration:      color.MustParse("red"),
		},
	}
	base := func(_ int, column string) string { return "<" + column + ">" }
	p := &AgeGradientTablePrinter{Table: NewTablePrinter(true, theme, base), Theme: theme}

	var buf bytes.Buffer
	p.Print(strings.NewReader("NAME  AGE\na     5m\n"), &buf)
	if got := buf.String(); !strings.Contains(got, "<5m>") {
		t.Errorf("single row: expected base filter to apply, got %q", got)
	}

	buf.Reset()
	p.Print(strings.NewReader("NAME  AGE\na     5m\nb     5m\n"), &buf)
	if got := buf.String(); strings.Count(got, "<5m>") != 2 {
		t.Errorf("equal ages: expected base filter to apply, got %q", got)
	}
}
