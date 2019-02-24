package writer

import (
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func Write(w io.Writer, styleSheet ast.StyleSheet) {
	writeStyleSheet(w, styleSheet)
}
