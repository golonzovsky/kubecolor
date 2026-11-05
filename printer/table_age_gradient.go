package printer

import (
	"bytes"
	"io"
	"log/slog"
	"math"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/kubecolor/kubecolor/scanner/tablescan"
)

var fallbackOldAgeColor = color.MustParse("gray")

// AgeGradientTablePrinter colors each duration column on a log-scale gradient
// from the first theme.data.durationFresh color (newest cell) to
// theme.data.duration (oldest cell). It buffers the whole input to learn each
// column's range, so it must not be used for streamed output like --watch.
type AgeGradientTablePrinter struct {
	Table *TablePrinter
	Theme *config.Theme
}

var _ Printer = &AgeGradientTablePrinter{}

func (p *AgeGradientTablePrinter) Print(r io.Reader, w io.Writer) {
	buf, err := io.ReadAll(r)
	if err != nil {
		slog.Error("Failed to read table output.", "error", err)
		return
	}
	ranges := scanAgeRanges(bytes.NewReader(buf))
	fresh, old := p.gradientEnds()

	table := *p.Table
	base := table.ColumnFilter
	table.ColumnFilter = func(i int, column string) string {
		if age, ok := stringutil.ParseHumanDuration(column); ok {
			if t, ok := ranges[i].fraction(age); ok {
				if c, ok := color.Interpolate(fresh, old, t); ok {
					return c.Render(column)
				}
			}
		}
		if base != nil {
			return base(i, column)
		}
		return column
	}
	table.Print(bytes.NewReader(buf), w)
}

func (p *AgeGradientTablePrinter) gradientEnds() (fresh, old color.Color) {
	if len(p.Theme.Data.DurationFresh) > 0 {
		fresh = p.Theme.Data.DurationFresh[0]
	}
	old = p.Theme.Data.Duration
	if _, ok := color.Interpolate(fresh, old, 0); !ok {
		old = fallbackOldAgeColor
	}
	return fresh, old
}

type ageRange struct {
	min, max time.Duration
}

func scanAgeRanges(r io.Reader) map[int]ageRange {
	ranges := make(map[int]ageRange)
	scanner := tablescan.NewScanner(r)
	for scanner.Scan() {
		for i, cell := range scanner.Cells() {
			age, ok := stringutil.ParseHumanDuration(cell.Trimmed)
			if !ok {
				continue
			}
			rng, found := ranges[i]
			if !found {
				ranges[i] = ageRange{min: age, max: age}
				continue
			}
			ranges[i] = ageRange{min: min(rng.min, age), max: max(rng.max, age)}
		}
	}
	return ranges
}

// fraction maps age onto a log scale from the column's newest (0) to oldest (1)
// value; sub-second ages clamp to 1s so the scale stays finite.
func (r ageRange) fraction(age time.Duration) (float64, bool) {
	lo := max(r.min, time.Second)
	hi := max(r.max, time.Second)
	if hi <= lo {
		return 0, false
	}
	age = min(max(age, lo), hi)
	return math.Log(float64(age)/float64(lo)) / math.Log(float64(hi)/float64(lo)), true
}
