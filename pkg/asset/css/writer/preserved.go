package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func writePreserved(w io.Writer, preserved ast.Preserved) {
	preserved.VisitPreserved(ast.PreservedVisitor{
		Ident: func(p ast.Ident) {
			fmt.Fprintf(w, p.Value)
		},

		AtKeyword: func(p ast.AtKeyword) {
			fmt.Fprintf(w, "@%s", p.Value)
		},

		Hash: func(p ast.Hash) {
			fmt.Fprintf(w, "#%s", p.Value)
		},

		String: func(p ast.String) {
			fmt.Fprintf(w, p.Value)
		},

		Url: func(p ast.Url) {
			fmt.Fprintf(w, "url(%s)", p.Value)
		},

		Number: func(p ast.Number) {
			if p.Integer {
				fmt.Fprintf(w, "%d", int(p.Value))
			} else {
				fmt.Fprintf(w, "%g", p.Value)
			}
		},

		Whitespace: func(ast.Whitespace) {
			fmt.Fprintf(w, " ")
		},

		Colon: func(ast.Colon) {
			fmt.Fprintf(w, ":")
		},

		Semicolon: func(ast.Semicolon) {
			fmt.Fprintf(w, ";")
		},
	})
}
