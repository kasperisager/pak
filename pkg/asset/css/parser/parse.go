package parser

import (
	"fmt"
	"net/url"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type SyntaxError struct {
	Offset  int
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

func Parse(tokens []token.Token) (*ast.StyleSheet, error) {
	offset, tokens, styleSheet, err := parseStyleSheet(0, tokens)

	if err != nil {
		return nil, err
	}

	if len(tokens) > 0 {
		return nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token",
		}
	}

	return styleSheet, nil
}

func parseStyleSheet(offset int, tokens []token.Token) (int, []token.Token, *ast.StyleSheet, error) {
	styleSheet := &ast.StyleSheet{}

	for {
		switch peek(tokens, 1).(type) {
		case token.Whitespace:
			offset, tokens = offset+1, tokens[1:]

		case token.CloseParen, token.CloseCurly, token.CloseSquare, nil:
			return offset, tokens, styleSheet, nil

		default:
			var (
				rule ast.Rule
				err  error
			)

			offset, tokens, rule, err = parseRule(offset, tokens)

			if err != nil {
				return offset, tokens, nil, err
			}

			styleSheet.Rules = append(styleSheet.Rules, rule)
		}
	}
}

func parseRule(offset int, tokens []token.Token) (int, []token.Token, ast.Rule, error) {
	switch t := peek(tokens, 1).(type) {
	case token.AtKeyword:
		switch t.Value {
		case "import":
			return parseImportRule(offset+1, tokens[1:])

		case "media":
			return parseMediaRule(offset+1, tokens[1:])

		case "keyframes":
			return parseKeyframesRule(offset+1, tokens[1:], "")
		case "-webkit-keyframes":
			return parseKeyframesRule(offset+1, tokens[1:], "-webkit-")

		case "supports":
			return parseSupportsRule(offset+1, tokens[1:])

		case "page":
			return parsePageRule(offset+1, tokens[1:])

		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token",
			}
		}

	default:
		return parseStyleRule(offset, tokens)
	}
}

func parseImportRule(offset int, tokens []token.Token) (int, []token.Token, *ast.ImportRule, error) {
	rule := &ast.ImportRule{}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Url:
		parsed, err := url.Parse(t.Value)

		if err != nil {
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: err.Error(),
			}
		}

		rule.URL = parsed
		offset, tokens = offset+1, tokens[1:]

	case token.String:
		parsed, err := url.Parse(t.Value)

		if err != nil {
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: err.Error(),
			}
		}

		rule.URL = parsed
		offset, tokens = offset+1, tokens[1:]

	case token.Function:
		if t.Value != "url" {
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: `unexpected function, expected "url()"`,
			}
		}

		offset, tokens = skipWhitespace(offset+1, tokens[1:])

		switch t := peek(tokens, 1).(type) {
		case token.String:
			parsed, err := url.Parse(t.Value)

			if err != nil {
				return offset, tokens, nil, SyntaxError{
					Offset:  offset,
					Message: err.Error(),
				}
			}

			rule.URL = parsed
			offset, tokens = offset+1, tokens[1:]

		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected string",
			}
		}

		offset, tokens = skipWhitespace(offset, tokens)

		switch peek(tokens, 1).(type) {
		case token.CloseParen:
			offset, tokens = offset+1, tokens[1:]

		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: `unexpected token, expected ")"`,
			}
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.Semicolon:
		return offset + 1, tokens[1:], rule, nil

	case nil:
		return offset, tokens, rule, nil

	default:
		offset, tokens, conditions, err := parseMediaQueryList(offset, tokens)

		if err != nil {
			return offset, tokens, rule, err
		}

		rule.Conditions = conditions

		offset, tokens = skipWhitespace(offset, tokens)

		switch peek(tokens, 1).(type) {
		case token.Semicolon:
			return offset + 1, tokens[1:], rule, nil

		case nil:
			return offset, tokens, rule, nil

		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: `unexpected token, expected ";"`,
			}
		}
	}
}

func parseMediaRule(offset int, tokens []token.Token) (int, []token.Token, *ast.MediaRule, error) {
	rule := &ast.MediaRule{}

	offset, tokens, conditions, err := parseMediaQueryList(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.Conditions = conditions

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.OpenCurly:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "{"`,
		}
	}

	offset, tokens, styleSheet, err := parseStyleSheet(offset, tokens)

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.StyleSheet = styleSheet

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseCurly:
		return offset + 1, tokens[1:], rule, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "}"`,
		}
	}
}

func parseStyleRule(offset int, tokens []token.Token) (int, []token.Token, *ast.StyleRule, error) {
	rule := &ast.StyleRule{}

	offset, tokens, selectors, err := parseSelectorList(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.Selectors = selectors

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.OpenCurly:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "{"`,
		}
	}

	offset, tokens, declarations, err := parseDeclarationList(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, rule, err
	}

	rule.Declarations = declarations

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseCurly:
		return offset + 1, tokens[1:], rule, nil

	default:
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "}"`,
		}
	}
}

func parseKeyframesRule(offset int, tokens []token.Token, prefix string) (int, []token.Token, *ast.KeyframesRule, error) {
	rule := &ast.KeyframesRule{Prefix: prefix}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		rule.Name = t.Value

	case token.String:
		rule.Name = t.Value

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected string or ident",
		}
	}

	offset, tokens = skipWhitespace(offset+1, tokens[1:])

	switch peek(tokens, 1).(type) {
	case token.OpenCurly:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "{"`,
		}
	}

	for {
		switch peek(tokens, 1).(type) {
		case token.Whitespace:
			offset, tokens = offset+1, tokens[1:]

		case token.CloseCurly:
			return offset + 1, tokens[1:], rule, nil

		default:
			var (
				block *ast.KeyframeBlock
				err   error
			)

			offset, tokens, block, err = parseKeyframeBlock(offset, tokens)

			if err != nil {
				return offset, tokens, rule, err
			}

			rule.Blocks = append(rule.Blocks, block)
		}
	}
}

func parseKeyframeBlock(offset int, tokens []token.Token) (int, []token.Token, *ast.KeyframeBlock, error) {
	block := &ast.KeyframeBlock{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		switch t.Value {
		case "from":
			block.Selector = 0
		case "to":
			block.Selector = 1
		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: `unexpected token, expected "from" or "to"`,
			}
		}

	case token.Percentage:
		block.Selector = t.Value

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "from", "to", or percentage`,
		}
	}

	offset, tokens = skipWhitespace(offset+1, tokens[1:])

	switch peek(tokens, 1).(type) {
	case token.OpenCurly:
		var (
			declarations []*ast.Declaration
			err          error
		)

		offset, tokens, declarations, err = parseDeclarationList(skipWhitespace(offset+1, tokens[1:]))

		if err != nil {
			return offset, tokens, nil, err
		}

		block.Declarations = declarations

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "{"`,
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseCurly:
		return offset + 1, tokens[1:], block, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "}"`,
		}
	}
}

func parseSupportsRule(offset int, tokens []token.Token) (int, []token.Token, *ast.SupportsRule, error) {
	rule := &ast.SupportsRule{}

	offset, tokens, condition, err := parseSupportsCondition(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.Condition = condition

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.OpenCurly:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "{"`,
		}
	}

	offset, tokens, styleSheet, err := parseStyleSheet(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.StyleSheet = styleSheet

	switch peek(tokens, 1).(type) {
	case token.CloseCurly:
		return offset + 1, tokens[1:], rule, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "}"`,
		}
	}
}

func parsePageRule(offset int, tokens []token.Token) (int, []token.Token, *ast.PageRule, error) {
	rule := &ast.PageRule{}

	offset, tokens, selectors, err := parsePageSelectorList(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.Selectors = selectors

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.OpenCurly:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "{"`,
		}
	}

	offset, tokens, components, err := parsePageComponentList(skipWhitespace(offset, tokens))

	if err != nil {
		return offset, tokens, nil, err
	}

	rule.Components = components

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseCurly:
		return offset + 1, tokens[1:], rule, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "}"`,
		}
	}
}

func parseDeclarationList(offset int, tokens []token.Token) (int, []token.Token, []*ast.Declaration, error) {
	var declarations []*ast.Declaration

	for {
		switch peek(tokens, 1).(type) {
		case token.Whitespace, token.Semicolon:
			offset, tokens = offset+1, tokens[1:]

		case token.CloseCurly:
			return offset, tokens, declarations, nil

		default:
			var (
				declaration *ast.Declaration
				err         error
			)

			offset, tokens, declaration, err = parseDeclaration(offset, tokens)

			if err != nil {
				return offset, tokens, nil, err
			}

			declarations = append(declarations, declaration)
		}
	}
}

func parseDeclaration(offset int, tokens []token.Token) (int, []token.Token, *ast.Declaration, error) {
	declaration := &ast.Declaration{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		declaration.Name = t.Value

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ident`,
		}
	}

	offset, tokens = skipWhitespace(offset+1, tokens[1:])

	switch peek(tokens, 1).(type) {
	case token.Colon:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ":"`,
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	for {
		switch peek(tokens, 1).(type) {
		case token.CloseParen, token.CloseCurly, token.CloseSquare, token.Semicolon, nil:
			return offset, tokens, declaration, nil

		default:
			var (
				component []token.Token
				err       error
			)

			offset, tokens, component, err = parseComponent(offset, tokens)

			if err != nil {
				return offset, tokens, nil, err
			}

			declaration.Value = append(declaration.Value, component...)
		}
	}
}

func parseComponent(offset int, tokens []token.Token) (int, []token.Token, []token.Token, error) {
	switch peek(tokens, 1).(type) {
	case token.OpenParen, token.Function:
		return parseParenBlock(offset+1, tokens[1:])

	case nil:
		return offset, tokens, nil, nil

	default:
		return offset + 1, tokens[1:], tokens[:1], nil
	}
}

func parseParenBlock(offset int, tokens []token.Token) (int, []token.Token, []token.Token, error) {
	var block []token.Token

	for {
		switch peek(tokens, 1).(type) {
		case token.CloseParen:
			return offset + 1, tokens[1:], block, nil

		case nil:
			return offset, tokens, block, SyntaxError{
				Offset:  offset,
				Message: `unexpected token, expected ")"`,
			}

		default:
			var (
				component []token.Token
				err       error
			)

			offset, tokens, component, err = parseComponent(offset, tokens)

			if err != nil {
				return offset, tokens, block, err
			}

			block = append(block, component...)
		}
	}
}

func parseSelectorList(offset int, tokens []token.Token) (int, []token.Token, []ast.Selector, error) {
	var selectors []ast.Selector

	for {
		switch peek(tokens, 1).(type) {
		case token.Whitespace, token.Comma:
			offset, tokens = offset+1, tokens[1:]

		case token.OpenCurly:
			return offset, tokens, selectors, nil

		default:
			var (
				selector ast.Selector
				err      error
			)

			offset, tokens, selector, err = parseSelector(offset, tokens)

			if err != nil {
				return offset, tokens, selectors, err
			}

			selectors = append(selectors, selector)
		}
	}
}

func parseSelector(offset int, tokens []token.Token) (int, []token.Token, ast.Selector, error) {
	var (
		left  ast.Selector
		right ast.Selector
		err   error
	)

	for {
		switch t := peek(tokens, 1).(type) {
		case token.Delim:
			switch t.Value {
			case '.':
				offset, tokens, right, err = parseClassSelector(offset+1, tokens[1:])

				if err != nil {
					return offset, tokens, left, err
				}

				left = combineSelectors(left, right)

			case '#':
				offset, tokens, right, err = parseIdSelector(offset+1, tokens[1:])

				if err != nil {
					return offset, tokens, left, err
				}

				left = combineSelectors(left, right)

			case '*':
				offset, tokens, right = offset+1, tokens[1:], &ast.TypeSelector{Name: "*"}

				left = combineSelectors(left, right)

			case '>', '~', '+':
				offset, tokens, left, err = parseRelativeSelector(offset, tokens, left)

				if err != nil {
					return offset, tokens, left, err
				}

			default:
				return offset, tokens, left, nil
			}

		case token.OpenSquare:
			offset, tokens, right, err = parseAttributeSelector(offset+1, tokens[1:])

			if err != nil {
				return offset, tokens, left, err
			}

			left = combineSelectors(left, right)

		case token.Ident:
			offset, tokens, right, err = parseTypeSelector(offset, tokens)

			if err != nil {
				return offset, tokens, left, err
			}

			left = combineSelectors(left, right)

		case token.Colon:
			offset, tokens, right, err = parsePseudoSelector(offset+1, tokens[1:])

			if err != nil {
				return offset, tokens, left, err
			}

			left = combineSelectors(left, right)

		case token.Whitespace:
			if len(tokens) > 1 && startsSelector(tokens[1]) {
				offset, tokens, left, err = parseRelativeSelector(offset, tokens, left)

				if err != nil {
					return offset, tokens, left, err
				}
			} else {
				offset, tokens = offset+1, tokens[1:]
			}

		default:
			return offset, tokens, left, nil
		}
	}
}

func startsSelector(t token.Token) bool {
	switch t := t.(type) {
	case token.Ident, token.Colon:
		return true

	case token.Delim:
		return t.Value == '.' || t.Value == '#' || t.Value == '*'
	}

	return false
}

func combineSelectors(left ast.Selector, right ast.Selector) ast.Selector {
	if left == nil {
		return right
	}

	return &ast.CompoundSelector{Left: left, Right: right}
}

func parseIdSelector(offset int, tokens []token.Token) (int, []token.Token, *ast.IdSelector, error) {
	selector := &ast.IdSelector{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Name = t.Value
		return offset + 1, tokens[1:], selector, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected id",
		}
	}
}

func parseClassSelector(offset int, tokens []token.Token) (int, []token.Token, *ast.ClassSelector, error) {
	selector := &ast.ClassSelector{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Name = t.Value
		return offset + 1, tokens[1:], selector, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected class",
		}
	}
}

func parseAttributeSelector(offset int, tokens []token.Token) (int, []token.Token, *ast.AttributeSelector, error) {
	selector := &ast.AttributeSelector{}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Name = t.Value
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected attribute name",
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.CloseSquare:
		return offset + 1, tokens[1:], selector, nil

	case token.Delim:
		switch t.Value {
		case '=':
			selector.Matcher = "="
			offset, tokens = offset+1, tokens[1:]

		case '~', '|', '^', '$', '*':
			offset, tokens = offset+1, tokens[1:]

			switch u := peek(tokens, 1).(type) {
			case token.Delim:
				if u.Value == '=' {
					switch t.Value {
					case '~':
						selector.Matcher = "~="
					case '|':
						selector.Matcher = "|="
					case '^':
						selector.Matcher = "^="
					case '$':
						selector.Matcher = "$="
					case '*':
						selector.Matcher = "*="
					}

					offset, tokens = offset+1, tokens[1:]
				}
			}
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.String:
		selector.Value = fmt.Sprintf("%c%s%[1]c", t.Mark, t.Value)

	case token.Ident:
		selector.Value = t.Value

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected attribute value",
		}
	}

	offset, tokens = skipWhitespace(offset+1, tokens[1:])

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		switch t.Value {
		case "i", "s":
			selector.Modifier = t.Value
			offset, tokens = offset+1, tokens[1:]
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseSquare:
		return offset + 1, tokens[1:], selector, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "]"`,
		}
	}
}

func parseTypeSelector(offset int, tokens []token.Token) (int, []token.Token, *ast.TypeSelector, error) {
	selector := &ast.TypeSelector{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Name = t.Value
		return offset + 1, tokens[1:], selector, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected type",
		}
	}
}

func parsePseudoSelector(offset int, tokens []token.Token) (int, []token.Token, *ast.PseudoSelector, error) {
	selector := &ast.PseudoSelector{Name: ":"}

	switch peek(tokens, 1).(type) {
	case token.Colon:
		selector.Name += ":"
		offset, tokens = offset+1, tokens[1:]
	}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Name += t.Value
		return offset + 1, tokens[1:], selector, nil

	case token.Function:
		selector.Name += t.Value
		offset, tokens = offset+1, tokens[1:]

		for {
			switch t := peek(tokens, 1).(type) {
			case token.CloseParen:
				return offset + 1, tokens[1:], selector, nil

			case nil:
				return offset, tokens, nil, SyntaxError{
					Offset:  offset,
					Message: `unexpected token, expected ")"`,
				}

			default:
				selector.Value = append(selector.Value, t)
				offset, tokens = offset+1, tokens[1:]
			}
		}

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected ident",
		}
	}
}

func parseRelativeSelector(offset int, tokens []token.Token, left ast.Selector) (int, []token.Token, *ast.RelativeSelector, error) {
	selector := &ast.RelativeSelector{Left: left}

	switch t := peek(tokens, 1).(type) {
	case token.Whitespace:
		selector.Combinator = ' '

	case token.Delim:
		switch t.Value {
		case '>', '~', '+':
			selector.Combinator = t.Value

		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected selector combinator",
			}
		}

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected selector combinator",
		}
	}

	var (
		right ast.Selector
		err   error
	)

	offset, tokens, right, err = parseSelector(skipWhitespace(offset+1, tokens[1:]))

	if err != nil {
		return offset, tokens, nil, err
	}

	selector.Right = right

	return offset, tokens, selector, nil
}

func parseMediaQueryList(offset int, tokens []token.Token) (int, []token.Token, []*ast.MediaQuery, error) {
	var mediaQueries []*ast.MediaQuery

	for {
		offset, tokens = skipWhitespace(offset, tokens)

		var (
			mediaQuery *ast.MediaQuery
			err        error
		)

		offset, tokens, mediaQuery, err = parseMediaQuery(offset, tokens)

		if err != nil {
			return offset, tokens, mediaQueries, err
		}

		mediaQueries = append(mediaQueries, mediaQuery)

		offset, tokens = skipWhitespace(offset, tokens)

		switch peek(tokens, 1).(type) {
		case token.Comma:
			offset, tokens = offset+1, tokens[1:]

		case token.OpenCurly, token.Semicolon, nil:
			return offset, tokens, mediaQueries, nil

		default:
			return offset, tokens, mediaQueries, SyntaxError{
				Offset:  offset,
				Message: "unexpected token",
			}
		}
	}
}

func parseMediaQuery(offset int, tokens []token.Token) (int, []token.Token, *ast.MediaQuery, error) {
	mediaQuery := &ast.MediaQuery{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		switch t.Value {
		case "not", "only":
			mediaQuery.Qualifier = t.Value
			offset, tokens = offset+1, tokens[1:]
		}

	case token.OpenParen:
		var (
			condition ast.MediaCondition
			err       error
		)

		offset, tokens, condition, err = parseMediaCondition(offset, tokens)

		if err != nil {
			return offset, tokens, nil, err
		}

		mediaQuery.Condition = condition

		return offset, tokens, mediaQuery, nil
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		switch t.Value {
		case "not":
		case "only":
		case "and":
		case "or":
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected media type",
			}

		default:
			mediaQuery.Type = t.Value
			offset, tokens = offset+1, tokens[1:]
		}

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token",
		}
	}

	return offset, tokens, mediaQuery, nil
}

func parseMediaCondition(offset int, tokens []token.Token) (int, []token.Token, ast.MediaCondition, error) {
	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		if t.Value == "not" {
			var (
				condition ast.MediaCondition
				err       error
			)

			offset, tokens, condition, err = parseMediaExpression(
				skipWhitespace(offset+1, tokens[1:]),
			)

			if err != nil {
				return offset, tokens, nil, err
			}

			condition = &ast.MediaNegation{Condition: condition}

			return offset, tokens, condition, nil
		}

		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token",
		}

	default:
		var (
			left ast.MediaCondition
			err  error
		)

		offset, tokens, left, err = parseMediaExpression(offset, tokens)

		if err != nil {
			return offset, tokens, nil, err
		}

		var operator string

		for {
			offset, tokens = skipWhitespace(offset, tokens)

			switch t := peek(tokens, 1).(type) {
			case token.Ident:
				switch t.Value {
				case "and", "or":
					if operator == "" || operator == t.Value {
						operator = t.Value
					} else {
						return offset, tokens, nil, SyntaxError{
							Offset:  offset,
							Message: `unexpected token, expected "` + operator + `"`,
						}
					}

					offset, tokens = offset+1, tokens[1:]

				default:
					return offset, tokens, nil, SyntaxError{
						Offset:  offset,
						Message: `unexpected token, expected "and" or "or"`,
					}
				}

				offset, tokens = skipWhitespace(offset, tokens)

				var right ast.MediaCondition

				offset, tokens, right, err = parseMediaExpression(offset, tokens)

				if err != nil {
					return offset, tokens, nil, err
				}

				left = &ast.MediaOperation{
					Operator: operator,
					Left:     left,
					Right:    right,
				}

			default:
				return offset, tokens, left, nil
			}
		}
	}
}

func parseMediaExpression(offset int, tokens []token.Token) (int, []token.Token, ast.MediaCondition, error) {
	if offset, tokens, feature, err := parseMediaFeature(offset, tokens); err == nil {
		return offset, tokens, feature, nil
	}

	switch peek(tokens, 1).(type) {
	case token.OpenParen:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "("`,
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	var (
		condition ast.MediaCondition
		err       error
	)

	offset, tokens, condition, err = parseMediaCondition(offset, tokens)

	if err != nil {
		return offset, tokens, condition, err
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseParen:
		return offset + 1, tokens[1:], condition, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ")"`,
		}
	}
}

func parseMediaFeature(offset int, tokens []token.Token) (int, []token.Token, *ast.MediaFeature, error) {
	feature := &ast.MediaFeature{}

	switch peek(tokens, 1).(type) {
	case token.OpenParen:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "("`,
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		feature.Name = t.Value
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected ident",
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.Colon:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ":"`,
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Number:
		feature.Value = &ast.MediaValuePlain{Value: t}

	case token.Dimension:
		feature.Value = &ast.MediaValuePlain{Value: t}

	case token.Ident:
		feature.Value = &ast.MediaValuePlain{Value: t}

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected number, dimension, or ident`,
		}
	}

	offset, tokens = skipWhitespace(offset+1, tokens[1:])

	switch peek(tokens, 1).(type) {
	case token.CloseParen:
		offset, tokens = offset+1, tokens[1:]
	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ")"`,
		}
	}

	return offset, tokens, feature, nil
}

func parseSupportsCondition(offset int, tokens []token.Token) (int, []token.Token, ast.SupportsCondition, error) {
	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		if t.Value == "not" {
			var (
				condition ast.SupportsCondition
				err       error
			)

			offset, tokens, condition, err = parseSupportsExpression(
				skipWhitespace(offset+1, tokens[1:]),
			)

			if err != nil {
				return offset, tokens, nil, err
			}

			condition = &ast.SupportsNegation{Condition: condition}

			return offset, tokens, condition, nil
		}

		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token",
		}

	default:
		var (
			left ast.SupportsCondition
			err  error
		)

		offset, tokens, left, err = parseSupportsExpression(offset, tokens)

		if err != nil {
			return offset, tokens, left, err
		}

		var operator string

		for {
			offset, tokens = skipWhitespace(offset, tokens)

			switch t := peek(tokens, 1).(type) {
			case token.Ident:
				switch t.Value {
				case "and", "or":
					if operator == "" || operator == t.Value {
						operator = t.Value
					} else {
						return offset, tokens, nil, SyntaxError{
							Offset:  offset,
							Message: `unexpected token, expected "` + operator + `"`,
						}
					}

					offset, tokens = offset+1, tokens[1:]

				default:
					return offset, tokens, nil, SyntaxError{
						Offset:  offset,
						Message: `unexpected token, expected "and" or "or"`,
					}
				}

				offset, tokens = skipWhitespace(offset, tokens)

				var right ast.SupportsCondition

				offset, tokens, right, err = parseSupportsExpression(offset, tokens)

				if err != nil {
					return offset, tokens, nil, err
				}

				left = &ast.SupportsOperation{
					Operator: operator,
					Left:     left,
					Right:    right,
				}

			default:
				return offset, tokens, left, nil
			}
		}
	}
}

func parseSupportsExpression(offset int, tokens []token.Token) (int, []token.Token, ast.SupportsCondition, error) {
	if offset, tokens, feature, err := parseSupportsFeature(offset, tokens); err == nil {
		return offset, tokens, feature, nil
	}

	switch peek(tokens, 1).(type) {
	case token.OpenParen:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "("`,
		}
	}

	var (
		condition ast.SupportsCondition
		err       error
	)

	offset, tokens, condition, err = parseSupportsCondition(offset, tokens)

	if err != nil {
		return offset, tokens, condition, err
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseParen:
		return offset + 1, tokens[1:], condition, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ")"`,
		}
	}
}

func parseSupportsFeature(offset int, tokens []token.Token) (int, []token.Token, *ast.SupportsFeature, error) {
	feature := &ast.SupportsFeature{}

	switch peek(tokens, 1).(type) {
	case token.OpenParen:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected "("`,
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	var (
		declaration *ast.Declaration
		err         error
	)

	offset, tokens, declaration, err = parseDeclaration(offset, tokens)

	if err != nil {
		return offset, tokens, nil, err
	}

	feature.Declaration = declaration

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseParen:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ")"`,
		}
	}

	return offset, tokens, feature, nil
}

func parsePageSelectorList(offset int, tokens []token.Token) (int, []token.Token, []*ast.PageSelector, error) {
	var selectors []*ast.PageSelector

	for {
		switch peek(tokens, 1).(type) {
		case token.Whitespace, token.Comma:
			offset, tokens = offset+1, tokens[1:]

		case token.OpenCurly:
			return offset, tokens, selectors, nil

		default:
			var (
				selector *ast.PageSelector
				err      error
			)

			offset, tokens, selector, err = parsePageSelector(offset, tokens)

			if err != nil {
				return offset, tokens, nil, err
			}

			selectors = append(selectors, selector)
		}
	}
}

func parsePageSelector(offset int, tokens []token.Token) (int, []token.Token, *ast.PageSelector, error) {
	selector := &ast.PageSelector{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Type = t.Value
		offset, tokens = offset+1, tokens[1:]
	}

	for {
		switch peek(tokens, 1).(type) {
		case token.Colon:
			switch t := peek(tokens, 2).(type) {
			case token.Ident:
				switch t.Value {
				case "left", "right", "first", "blank":
					selector.Classes = append(selector.Classes, ":"+t.Value)
					offset, tokens = offset+2, tokens[2:]

				default:
					return offset, tokens, nil, SyntaxError{
						Offset:  offset,
						Message: `unexpected token, expected page selector`,
					}
				}

			default:
				return offset, tokens, nil, SyntaxError{
					Offset:  offset,
					Message: `unexpected token, expected page selector`,
				}
			}

		default:
			return offset, tokens, selector, nil
		}
	}
}

func parsePageComponentList(offset int, tokens []token.Token) (int, []token.Token, []ast.PageComponent, error) {
	var pageComponents []ast.PageComponent

	for {
		switch peek(tokens, 1).(type) {
		case token.Whitespace, token.Semicolon:
			offset, tokens = offset+1, tokens[1:]

		case token.CloseCurly:
			return offset, tokens, pageComponents, nil

		default:
			var (
				pageComponent ast.PageComponent
				err           error
			)

			offset, tokens, pageComponent, err = parsePageComponent(offset, tokens)

			if err != nil {
				return offset, tokens, nil, err
			}

			pageComponents = append(pageComponents, pageComponent)
		}
	}
}

func parsePageComponent(offset int, tokens []token.Token) (int, []token.Token, ast.PageComponent, error) {
	switch peek(tokens, 1).(type) {
	case token.Ident:
		var (
			declaration *ast.Declaration
			err         error
		)

		offset, tokens, declaration, err = parseDeclaration(offset, tokens)

		if err != nil {
			return offset, tokens, nil, err
		}

		return offset, tokens, &ast.PageDeclaration{Declaration: declaration}, nil

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: `unexpected token, expected ident`,
		}
	}
}

func peek(tokens []token.Token, n int) token.Token {
	if len(tokens) < n {
		return nil
	}

	return tokens[n-1]
}

func skipWhitespace(offset int, tokens []token.Token) (int, []token.Token) {
	if _, ok := peek(tokens, 1).(token.Whitespace); ok {
		return offset + 1, tokens[1:]
	}

	return offset, tokens
}
