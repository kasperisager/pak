package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func writeBlock(w io.Writer, block ast.Block) {
	fmt.Fprintf(w, "{")

	for _, preserved := range block.Value {
		writePreserved(w, preserved)
	}

	fmt.Fprintf(w, "}")
}
