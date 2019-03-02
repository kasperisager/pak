package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func writeRule(w io.Writer, rule ast.Rule) {
	switch rule := rule.(type) {
	case ast.AtRule:
		fmt.Fprintf(w, "@%s", rule.Name)

		if len(rule.Prelude) > 0 {
			fmt.Fprintf(w, " ")
		}

		for _, token := range rule.Prelude {
			writeToken(w, token)
		}

		if rule.Value != nil {
			writeBlock(w, *rule.Value)
		} else {
			fmt.Fprintf(w, ";")
		}

	case ast.QualifiedRule:
		for _, token := range rule.Prelude {
			writeToken(w, token)
		}

		writeBlock(w, rule.Value)
	}
}
