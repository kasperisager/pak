package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/css/token"
)

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
		if t.Integer {
			fmt.Fprintf(w, "%d", int(t.Value))
		} else {
			fmt.Fprintf(w, "%g", t.Value)
		}

	case token.Percentage:
		fmt.Fprintf(w, "%f%%", t.Value)

	case token.Dimension:
		if t.Integer {
			fmt.Fprintf(w, "%d%s", int(t.Value), t.Unit)
		} else {
			fmt.Fprintf(w, "%g%s", t.Value, t.Unit)
		}

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
