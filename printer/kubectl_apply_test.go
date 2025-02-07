package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golonzovsky/kubecolor/testutil"
)

func Test_ApplyPrinter_Print(t *testing.T) {
	tests := []struct {
		name         string
		tablePrinter *TablePrinter
		input        string
		expected     string
	}{
		{
			name: "created",
			input: testutil.NewHereDoc(`
				deployment.apps/foo created`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [32mcreated[0m
			`),
		},
		{
			name: "configured",
			input: testutil.NewHereDoc(`
				deployment.apps/foo configured`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [33mconfigured[0m
			`),
		},
		{
			name: "unchanged",
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [35munchanged[0m
			`),
		},
		{
			name: "dry run",
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged (dry run)`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [35munchanged[0m [36m(dry run)[0m
			`),
		},
		{
			name: "something else. This likely won't happen but fallbacks here just in case.",
			input: testutil.NewHereDoc(`
				deployment.apps/foo bar`),
			expected: testutil.NewHereDoc(`
				[32mdeployment.apps/foo bar[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := ApplyPrinter{}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
