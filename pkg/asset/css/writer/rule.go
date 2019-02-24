package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func writeRule(w io.Writer, rule ast.Rule) {
	rule.VisitRule(ast.RuleVisitor{
		AtRule: func(rule ast.AtRule) {
			fmt.Fprintf(w, "@%s", rule.Name)

			if len(rule.Prelude) > 0 {
				fmt.Fprintf(w, " ")
			}

			for _, preserved := range rule.Prelude {
				writePreserved(w, preserved)
			}

			if rule.Value != nil {
				writeBlock(w, *rule.Value)
			} else {
				fmt.Fprintf(w, ";")
			}
		},

		QualifiedRule: func(rule ast.QualifiedRule) {
			for _, preserved := range rule.Prelude {
				writePreserved(w, preserved)
			}

			writeBlock(w, rule.Value)
		},
	})
}
