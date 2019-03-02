package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func writeBlock(w io.Writer, block ast.Block) {
	fmt.Fprintf(w, "{")

	for _, token := range block.Value {
		writeToken(w, token)
	}

	fmt.Fprintf(w, "}")
}
