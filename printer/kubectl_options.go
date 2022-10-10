package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

type OptionsPrinter struct{}

func (op *OptionsPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	isFirstLine := true
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		if isFirstLine {
			fmt.Fprintf(w, "%s\n", color.Apply(line, op.firstLineColor()))
			isFirstLine = false
			continue
		}

		indentCnt := findIndent(line)
		indent := toSpaces(indentCnt)
		trimmedLine := strings.TrimLeft(line, " ")

		splitted := strings.SplitN(trimmedLine, ": ", 2)
		key, val := splitted[0], splitted[1]

		fmt.Fprintf(w, "%s%s: %s\n", indent, color.Apply(key, getColorByKeyIndent(0, 2)), color.Apply(val, getColorByValueType(val)))
	}
}

func (op *OptionsPrinter) firstLineColor() color.Color {
	return StringColor
}
