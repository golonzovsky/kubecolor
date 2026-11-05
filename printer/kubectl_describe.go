package printer

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// DescribePrinter is used on "kubectl describe" output
type DescribePrinter struct {
	TablePrinter *TablePrinter

	tableBytes *bytes.Buffer
}

// ensures it implements the interface
var _ Printer = &DescribePrinter{}

var onlyValuePathToColor = regexp.MustCompile(`^(Labels|Annotations|(Init )?Containers/[^/]+/Environment Variables from)/.+`)

// Print implements [Printer.Print]
func (p *DescribePrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	const basicIndentWidth = 2 // according to kubectl describe format
	for scanner.Scan() {
		line := scanner.Line()

		if onlyValuePathToColor.MatchString(scanner.Path().String()) { // if line path matches label or annotation or env from
			line.Value = bytes.Join([][]byte{line.Key, line.Value}, line.Spacing)
			line.Key = nil
			line.Spacing = nil
		} else if bytesutil.CountColumns(line.Value, " \t") >= 3 { // when there are multiple columns, treat is as table format
			if p.tableBytes == nil {
				p.tableBytes = &bytes.Buffer{}
			}
			fmt.Fprintln(p.tableBytes, line.String())
			continue
		} else if p.tableBytes != nil {
			p.TablePrinter.Print(p.tableBytes, w)
			p.tableBytes = nil
		}

		fmt.Fprintf(w, "%s", line.Indent)
		if len(line.Key) > 0 {
			keyColor := ColorDataKey(line.KeyIndent(), basicIndentWidth, p.TablePrinter.Theme.Describe.Key)
			key := string(line.Key)
			if withoutColon, ok := strings.CutSuffix(key, ":"); ok {
				fmt.Fprint(w, keyColor.Render(withoutColon), ":")
			} else {
				fmt.Fprint(w, keyColor.Render(key))
			}
		}
		fmt.Fprintf(w, "%s", line.Spacing)
		if len(line.Value) > 0 {
			fmt.Fprint(w, p.colorizeValue(scanner.Path(), string(line.Value)))
		}
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print describe output.", "error", err)
	}

	if p.tableBytes != nil {
		p.TablePrinter.Print(p.tableBytes, w)
		p.tableBytes = nil
	}
}

func (p *DescribePrinter) colorize(path describe.Path, value string) string {
	value = strings.TrimSpace(value)
	pathStr := path.String()
	if col, ok := p.colorizeStatus(value, pathStr); ok {
		return col
	}
	if col, ok := p.colorizeArgs(value, pathStr); ok {
		return col
	}
	return ColorDataValue(value, p.TablePrinter.Theme).Render(value)
}

var describePathsToColor = []*regexp.Regexp{
	regexp.MustCompile(`^Status$`),
	regexp.MustCompile(`^(Init )?Containers/[^/]*/State(/Reason)?$`),
	regexp.MustCompile(`^Containers/[^/]*/Last State(/Reason)?$`),
}

func (p *DescribePrinter) colorizeStatus(value, pathStr string) (string, bool) {
	if !matchesAnyRegex(pathStr, describePathsToColor) {
		return value, false
	}
	col, ok := ColorStatus(value, p.TablePrinter.Theme)
	if !ok {
		return value, false
	}
	return col, true
}

var describePathsToArgs = []*regexp.Regexp{
	regexp.MustCompile(`^(Init )?Containers/[^/]*/Args$`),
}

func (p *DescribePrinter) colorizeArgs(value, pathStr string) (string, bool) {
	if !matchesAnyRegex(pathStr, describePathsToArgs) {
		return value, false
	}
	if !strings.HasPrefix(value, "-") {
		return value, false
	}
	return p.pairKey(value), true // same as the "--my-flag=" part of "--my-flag=123"
}

func matchesAnyRegex(s string, regexes []*regexp.Regexp) bool {
	for _, r := range regexes {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

var describePathsTolerations = []*regexp.Regexp{
	regexp.MustCompile(`(^|/)Tolerations$`),
	regexp.MustCompile(`(^|/)Taints$`),
}

func (p *DescribePrinter) colorizeValue(path describe.Path, val string) string {
	if matchesAnyRegex(path.String(), describePathsTolerations) {
		return mapTokens(val, func(i int, tok string) string { return p.colorizeToleration(path, i, tok) })
	}
	if k, v, ok := strings.Cut(val, ": "); ok { // annotations and env
		return p.pairKey(k+": ") + p.colorizePairsOrValue(path, v)
	}
	return p.colorizePairsOrValue(path, val)
}

// colorizePairsOrValue colors "key=value" pairs (labels, probes) with the key
// left uncolored, and any other value by its type.
func (p *DescribePrinter) colorizePairsOrValue(path describe.Path, val string) string {
	if !strings.Contains(val, "=") {
		return p.colorize(path, val)
	}
	return mapTokens(val, func(_ int, tok string) string {
		if k, v, ok := strings.Cut(tok, "="); ok {
			return p.pairKey(k+"=") + p.colorize(path, v)
		}
		return p.colorize(path, tok)
	})
}

// colorizeToleration colors one token of "key[=value][:effect] [op=Op] [for Ns]"
func (p *DescribePrinter) colorizeToleration(path describe.Path, i int, tok string) string {
	if i == 0 {
		keyValue, effect, hasEffect := strings.Cut(tok, ":")
		key, value, hasValue := strings.Cut(keyValue, "=")
		var sb strings.Builder
		switch {
		case hasValue:
			sb.WriteString(p.pairKey(key + "="))
			sb.WriteString(p.colorize(path, value))
			if hasEffect {
				sb.WriteString(p.pairKey(":"))
			}
		case hasEffect:
			sb.WriteString(p.pairKey(key + ":"))
		default:
			sb.WriteString(p.pairKey(key))
		}
		if hasEffect {
			sb.WriteString(p.colorize(path, effect))
		}
		return sb.String()
	}
	if k, v, ok := strings.Cut(tok, "="); ok {
		return p.pairKey(k+"=") + p.colorize(path, v)
	}
	if _, ok := stringutil.ParseHumanDuration(tok); ok {
		return p.colorize(path, tok)
	}
	return p.pairKey(tok)
}

func (p *DescribePrinter) pairKey(s string) string {
	if s == "" {
		return ""
	}
	return p.TablePrinter.Theme.Describe.PairKey.Render(s)
}

// mapTokens applies fn to each token of s split on spaces, tabs and commas,
// keeping the separators as they were
func mapTokens(s string, fn func(i int, tok string) string) string {
	var sb strings.Builder
	start, n := 0, 0
	for i := 0; i <= len(s); i++ {
		if i < len(s) && s[i] != ' ' && s[i] != '\t' && s[i] != ',' {
			continue
		}
		if tok := s[start:i]; tok != "" {
			sb.WriteString(fn(n, tok))
			n++
		}
		if i < len(s) {
			sb.WriteByte(s[i])
		}
		start = i + 1
	}
	return sb.String()
}
