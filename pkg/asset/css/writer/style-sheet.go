package writer

import (
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func writeStyleSheet(w io.Writer, styleSheet ast.StyleSheet) {
	for _, rule := range styleSheet.Rules {
		writeRule(w, rule)
	}
}
