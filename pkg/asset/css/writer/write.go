package writer

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func Write(w io.Writer, styleSheet ast.StyleSheet) {
	writeStyleSheet(w, styleSheet)
}

func writeStyleSheet(w io.Writer, styleSheet ast.StyleSheet) {
	for _, rule := range styleSheet.Rules {
		writeRule(w, rule)
	}
}

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

func writeBlock(w io.Writer, block ast.Block) {
	fmt.Fprintf(w, "{")

	for _, token := range block.Value {
		writeToken(w, token)
	}

	fmt.Fprintf(w, "}")
}

func writeToken(w io.Writer, t token.Token) {
	switch t := t.(type) {
	case token.Ident:
		fmt.Fprintf(w, t.Value)

	case token.Function:
		fmt.Fprintf(w, "%s(", t.Value)

	case token.AtKeyword:
		fmt.Fprintf(w, "@%s", t.Value)

	case token.Hash:
		fmt.Fprintf(w, "#%s", t.Value)

	case token.String:
		fmt.Fprintf(w, "\"%s\"", t.Value)

	case token.Url:
		fmt.Fprintf(w, "url(%s)", t.Value)

	case token.Number:
		fmt.Fprintf(w, "%s", strconv.FormatFloat(t.Value, 'f', -1, 64))

	case token.Percentage:
		fmt.Fprintf(w, "%s%%", strconv.FormatFloat(t.Value, 'f', -1, 64))

	case token.Dimension:
		fmt.Fprintf(w, "%s%s", strconv.FormatFloat(t.Value, 'f', -1, 64), t.Unit)

	case token.Whitespace:
		fmt.Fprintf(w, " ")

	case token.Colon:
		fmt.Fprintf(w, ":")

	case token.Semicolon:
		fmt.Fprintf(w, ";")

	case token.Comma:
		fmt.Fprintf(w, ",")

	case token.OpenSquare:
		fmt.Fprintf(w, "[")

	case token.CloseSquare:
		fmt.Fprintf(w, "]")

	case token.OpenParen:
		fmt.Fprintf(w, "(")

	case token.CloseParen:
		fmt.Fprintf(w, ")")

	case token.OpenCurly:
		fmt.Fprintf(w, "{")

	case token.CloseCurly:
		fmt.Fprintf(w, "}")
	}
}
