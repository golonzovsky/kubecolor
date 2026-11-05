package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/kubectl"
)

// Age-based coloring must apply only to "kubectl get" tables, not to other
// table subcommands like "kubectl events" (whose age-ish columns, e.g.
// "43s (x150 over 21h)", would color inconsistently).
func Test_KubectlOutputColoredPrinter_ageColorScopedToGet(t *testing.T) {
	theme := &config.Theme{
		Table: config.ThemeTable{
			Header:  color.MustParse("bold"),
			Columns: color.MustParseSlice("blue"), // column-cycle color
		},
		Data: config.ThemeData{
			DurationFresh: color.MustParseSlice("red"), // distinctive fresh color
		},
	}

	input := "AGE\n1m\n"
	render := func(sub kubectl.Subcommand) string {
		p := &KubectlOutputColoredPrinter{
			SubcommandInfo:    &kubectl.SubcommandInfo{Subcommand: sub},
			ObjFreshThreshold: config.MustParseDurationSlice("5m"),
			Theme:             theme,
		}
		var buf bytes.Buffer
		p.Print(strings.NewReader(input), &buf)
		return buf.String()
	}

	const freshCode = "\x1b[31m" // red

	get := render(kubectl.Get)
	if !strings.Contains(get, freshCode) {
		t.Errorf("get: expected age to be fresh-colored, got %q", get)
	}

	events := render(kubectl.Events)
	if strings.Contains(events, freshCode) {
		t.Errorf("events: expected age NOT fresh-colored, got %q", events)
	}
}

// The age gradient needs the whole table up front, so it must only wrap
// non-streaming "kubectl get" tables; --watch keeps the threshold coloring.
func Test_KubectlOutputColoredPrinter_ageGradientSkipsWatch(t *testing.T) {
	theme := &config.Theme{
		Data: config.ThemeData{
			DurationFresh: color.MustParseSlice("#0000ff"),
			Duration:      color.MustParse("#ff0000"),
		},
	}
	input := "AGE\n1m\n1d\n"
	render := func(watch bool) string {
		p := &KubectlOutputColoredPrinter{
			SubcommandInfo: &kubectl.SubcommandInfo{Subcommand: kubectl.Get, Watch: watch},
			AgeGradient:    true,
			Theme:          theme,
		}
		var buf bytes.Buffer
		p.Print(strings.NewReader(input), &buf)
		return buf.String()
	}
	fresh := color.MustParse("#0000ff").Render("1m")
	old := color.MustParse("#ff0000").Render("1d")
	if got := render(false); !strings.Contains(got, fresh) || !strings.Contains(got, old) {
		t.Errorf("get: expected gradient ends, got %q", got)
	}
	// no thresholds -> every age gets theme.data.duration
	if got := render(true); strings.Contains(got, fresh) || !strings.Contains(got, color.MustParse("#ff0000").Render("1m")) {
		t.Errorf("get --watch: expected threshold coloring, got %q", got)
	}
}
