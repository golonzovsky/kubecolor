package printer

import (
	"fmt"
	"io"
	"log"

	"github.com/andreazorzetto/yh/highlight"
)

type YamlPrinter struct {
}

func (yp *YamlPrinter) Print(r io.Reader, w io.Writer) {
	h, err := highlight.Highlight(r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, h)
}
