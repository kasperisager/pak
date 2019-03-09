package parser

import (
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type SyntaxError struct {
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

func Parse(tokens []token.Token) (ast.StyleSheet, error) {
	_, styleSheet, err := parseStyleSheet(tokens)

	return styleSheet, err
}

func parseStyleSheet(tokens []token.Token) ([]token.Token, ast.StyleSheet, error) {
	styleSheet := ast.StyleSheet{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Whitespace:
			tokens = tokens[1:]

		default:
			var (
				rule ast.Rule
				err  error
			)

			tokens, rule, err = parseRule(tokens)

			if err != nil {
				return tokens, styleSheet, err
			}

			styleSheet.Rules = append(styleSheet.Rules, rule)
		}
	}

	return tokens, styleSheet, nil
}

func parseRule(tokens []token.Token) ([]token.Token, ast.Rule, error) {
	switch t := tokens[0].(type) {
	case token.AtKeyword:
		switch t.Value {
		case "import":
			return parseImportRule(tokens[1:])
		}

	default:
		return parseStyleRule(tokens)
	}

	return tokens, nil, SyntaxError{
		Message: "unexpected token",
	}
}

func parseImportRule(tokens []token.Token) ([]token.Token, ast.ImportRule, error) {
	rule := ast.ImportRule{}

	tokens = skipWhitespace(tokens)

	if len(tokens) == 0 {
		return tokens, rule, SyntaxError{
			Message: "unexpected end of file, expected url",
		}
	}

	switch t := tokens[0].(type) {
	case token.Url:
		rule.Url = t.Value
		tokens = tokens[1:]

	case token.String:
		rule.Url = t.Value
		tokens = tokens[1:]

	case token.Function:
		if t.Value != "url" {
			return tokens, rule, SyntaxError{
				Message: "unexpected function, expected url()",
			}
		}

		tokens = skipWhitespace(tokens[1:])

		if len(tokens) == 0 {
			return tokens, rule, SyntaxError{
				Message: "unexpected end of file",
			}
		}

		if t, ok := tokens[0].(token.String); ok {
			rule.Url = t.Value
			tokens = tokens[1:]
		} else {
			return tokens, rule, SyntaxError{
				Message: "unexpected token, expected string",
			}
		}

		if len(tokens) == 0 {
			return tokens, rule, SyntaxError{
				Message: "unexpected end of file",
			}
		}

		if _, ok := tokens[0].(token.CloseParen); ok {
			tokens = tokens[1:]
		} else {
			return tokens, rule, SyntaxError{
				Message: "unexpected token, expected closing paren",
			}
		}
	}

	tokens = skipWhitespace(tokens)

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.Semicolon); ok {
			tokens = tokens[1:]
		} else {
			return tokens, rule, SyntaxError{
				Message: "unexpected token, expected semicolon",
			}
		}
	}

	return tokens, rule, nil
}


func parseStyleRule(tokens []token.Token) ([]token.Token, ast.StyleRule, error) {
	rule := ast.StyleRule{}

	if len(tokens) == 0 {
		return tokens, rule, SyntaxError{
			Message: "unexpected end of file, expected selector",
		}
	}

	tokens, selectors, err := parseSelectorList(tokens)

	if err != nil {
		return tokens, rule, err
	}

	rule.Selectors = selectors

	tokens = skipWhitespace(tokens)

	if len(tokens) == 0 {
		return tokens, rule, SyntaxError{
			Message: "unexpected end of file, expected selector",
		}
	}

	if _, ok := tokens[0].(token.OpenCurly); ok {
		tokens = tokens[1:]
	} else {
		return tokens, rule, SyntaxError{
			Message: "unexpected token, expected opening curly",
		}
	}

	tokens, declarations, err := parseDeclarationList(tokens)

	if err != nil {
		return tokens, rule, err
	}

	rule.Declarations = declarations

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.CloseCurly); ok {
			tokens = tokens[1:]
		} else {
			return tokens, rule, SyntaxError{
				Message: "unexpected token, expected closing curly",
			}
		}
	}

	return tokens, rule, nil
}

func parseSelectorList(tokens []token.Token) ([]token.Token, []ast.Selector, error) {
	selectors := []ast.Selector{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Whitespace, token.Comma:
			tokens = tokens[1:]

		case token.OpenCurly:
			return tokens, selectors, nil

		default:
			var (
				selector ast.Selector
				err error
			)

			tokens, selector, err = parseSelector(tokens)

			if err != nil {
				return tokens, selectors, err
			}

			selectors = append(selectors, selector)
		}
	}

	return tokens, selectors, SyntaxError{
		Message: "unexpected end of file",
	}
}

func parseSelector(tokens []token.Token) ([]token.Token, ast.Selector, error) {
	var (
		left ast.Selector
		right ast.Selector
		err error
	)

	for len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Comma, token.OpenCurly:
			return tokens, left, nil

		case token.Delim:
			switch t.Value {
			case '.':
				tokens, right, err = parseClassSelector(tokens[1:])

				if err != nil {
					return tokens, left, err
				}

				left = combineSelectors(left, right)

			case '#':
				tokens, right, err = parseIdSelector(tokens[1:])

				if err != nil {
					return tokens, left, err
				}

				left = combineSelectors(left, right)

			case '>', '~', '+':
				tokens, left, err = parseRelativeSelector(tokens, left)

				if err != nil {
					return tokens, left, err
				}

			default:
				return tokens, left, SyntaxError{
					Message: "unexpected token",
				}
			}

		case token.Whitespace:
			if len(tokens) > 1 && startsSelector(tokens[1]) {
				tokens, left, err = parseRelativeSelector(tokens, left)

				if err != nil {
					return tokens, left, err
				}
			} else {
				tokens = tokens[1:]
			}

		default:
			return tokens, left, SyntaxError{
				Message: "unexpected token",
			}
		}
	}

	return tokens, left, SyntaxError{
		Message: "unexpected end of file",
	}
}

func startsSelector(t token.Token) bool {
	switch t := t.(type) {
	case token.Delim:
		return t.Value == '.' || t.Value == '#'
	}

	return false
}

func combineSelectors(left ast.Selector, right ast.Selector) ast.Selector {
	if left == nil {
		return right
	}

	return ast.CompoundSelector{Left: left, Right: right}
}

func parseIdSelector(tokens []token.Token) ([]token.Token, ast.IdSelector, error) {
	selector := ast.IdSelector{}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Ident:
			selector.Name = t.Value
			return tokens[1:], selector, nil

		default:
			return tokens, selector, SyntaxError{
				Message: "unexpected token, expected id",
			}
		}
	}

	return tokens, selector, SyntaxError{
		Message: "unexpected end of file",
	}
}

func parseClassSelector(tokens []token.Token) ([]token.Token, ast.ClassSelector, error) {
	selector := ast.ClassSelector{}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Ident:
			selector.Name = t.Value
			return tokens[1:], selector, nil

		default:
			return tokens, selector, SyntaxError{
				Message: "unexpected token, expected class",
			}
		}
	}

	return tokens, selector, SyntaxError{
		Message: "unexpected end of file",
	}
}

func parseRelativeSelector(tokens []token.Token, left ast.Selector) ([]token.Token, ast.RelativeSelector, error) {
	selector := ast.RelativeSelector{Left: left}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Whitespace:
			selector.Combinator = ast.Descendant

		case token.Delim:
			switch t.Value {
			case '>':
				selector.Combinator = ast.DirectDescendant
			case '~':
				selector.Combinator = ast.Sibling
			case '+':
				selector.Combinator = ast.DirectSibling
			default:
				return tokens, selector, SyntaxError{
					Message: "unexpected token, expected selector combinator",
				}
			}

		default:
			return tokens, selector, SyntaxError{
				Message: "unexpected token, expected selector combinator",
			}
		}

		tokens = skipWhitespace(tokens[1:])
	}

	var (
		right ast.Selector
		err error
	)

	tokens, right, err = parseSelector(tokens)

	if err != nil {
		return tokens, selector, err
	}

	selector.Right = right

	return tokens, selector, nil
}

func parseDeclarationList(tokens []token.Token) ([]token.Token, []ast.Declaration, error) {
	declarations := []ast.Declaration{}

	for len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Whitespace, token.Semicolon:
			tokens = tokens[1:]

		case token.CloseCurly:
			return tokens, declarations, nil

		case token.Ident:
			var (
				declaration ast.Declaration
				err error
			)

			tokens, declaration, err = parseDeclaration(tokens[1:], t.Value)

			if err != nil {
				return tokens, declarations, err
			}

			declarations = append(declarations, declaration)

		default:
			return tokens, declarations, SyntaxError{
				Message: "unexpected token",
			}
		}
	}

	return tokens, declarations, SyntaxError{
		Message: "unexpected end of file",
	}
}

func parseDeclaration(tokens []token.Token, name string) ([]token.Token, ast.Declaration, error) {
	declaration := ast.Declaration{Name: name}

	tokens = skipWhitespace(tokens)

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.Colon); ok {
			tokens = tokens[1:]
		} else {
			return tokens, declaration, SyntaxError{
				Message: "unexpected token, expected colon",
			}
		}
	}

	tokens = skipWhitespace(tokens)

	for len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.CloseCurly:
			return tokens, declaration, nil

		case token.Semicolon:
			return tokens[1:], declaration, nil

		default:
			declaration.Value = append(declaration.Value, t)
			tokens = tokens[1:]
		}
	}

	return tokens, declaration, SyntaxError{
		Message: "unexpected end of file",
	}
}

func skipWhitespace(tokens []token.Token) []token.Token {
	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.Whitespace); ok {
			return tokens[1:]
		}
	}

	return tokens
}
