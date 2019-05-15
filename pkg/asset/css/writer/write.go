package writer

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/token"
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
	case ast.ImportRule:
		fmt.Fprintf(w, "@import \"%s\"", rule.URL.String())

		for i, mediaQuery := range rule.Conditions {
			if i == 0 {
				fmt.Fprintf(w, " ")
			} else {
				fmt.Fprintf(w, ",")
			}

			writeMediaQuery(w, mediaQuery)
		}

		fmt.Fprintf(w, ";")

	case ast.StyleRule:
		for i, selector := range rule.Selectors {
			if i != 0 {
				fmt.Fprintf(w, ",")
			}

			writeSelector(w, selector)
		}

		fmt.Fprintf(w, "{")

		for i, declaration := range rule.Declarations {
			if i != 0 {
				fmt.Fprintf(w, ";")
			}

			writeDeclaration(w, declaration)
		}

		fmt.Fprintf(w, "}")

	case ast.MediaRule:
		fmt.Fprintf(w, "@media ")

		for i, mediaQuery := range rule.Conditions {
			if i != 0 {
				fmt.Fprintf(w, ",")
			}

			writeMediaQuery(w, mediaQuery)
		}

		fmt.Fprintf(w, "{")

		writeStyleSheet(w, rule.StyleSheet)

		fmt.Fprintf(w, "}")
	}
}

func writeSelector(w io.Writer, selector ast.Selector) {
	switch selector := selector.(type) {
	case ast.IdSelector:
		fmt.Fprintf(w, "#%s", selector.Name)

	case ast.ClassSelector:
		fmt.Fprintf(w, ".%s", selector.Name)

	case ast.AttributeSelector:
		fmt.Fprintf(w, "[")

		if selector.Namespace != nil {
			fmt.Fprintf(w, "%s|", *selector.Namespace)
		}

		fmt.Fprintf(w, "%s%s%s", selector.Name, selector.Matcher, selector.Value)

		if selector.Modifier != "" {
			fmt.Fprintf(w, " %s", selector.Modifier)
		}

		fmt.Fprintf(w, "]")

	case ast.TypeSelector:
		fmt.Fprintf(w, "%s", selector.Name)

	case ast.CompoundSelector:
		writeSelector(w, selector.Left)
		writeSelector(w, selector.Right)

	case ast.RelativeSelector:
		writeSelector(w, selector.Left)
		fmt.Fprintf(w, "%c", selector.Combinator)
		writeSelector(w, selector.Right)
	}
}

func writeMediaQuery(w io.Writer, mediaQuery ast.MediaQuery) {
	var parts []string

	if mediaQuery.Qualifier != "" {
		parts = append(parts, mediaQuery.Qualifier)
	}

	if mediaQuery.Type != "" {
		parts = append(parts, mediaQuery.Type)
	}

	fmt.Fprintf(w, "%s", strings.Join(parts, " "))
}

func writeDeclaration(w io.Writer, declaration ast.Declaration) {
	fmt.Fprintf(w, "%s:", declaration.Name)

	for _, t := range declaration.Value {
		writeToken(w, t)
	}
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
		fmt.Fprintf(w, "%c%s%[1]c", t.Mark, t.Value)

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
