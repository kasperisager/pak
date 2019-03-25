package parser

import (
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

func Parse(tokens []token.Token) (ast.StyleSheet, error) {
	offset, tokens, styleSheet, err := parseStyleSheet(0, tokens)

	if len(tokens) > 0 {
		return styleSheet, SyntaxError{
			Offset:  offset,
			Message: "unexpected token",
		}
	}

	return styleSheet, err
}

func parseStyleSheet(offset int, tokens []token.Token) (int, []token.Token, ast.StyleSheet, error) {
	styleSheet := ast.StyleSheet{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Whitespace:
			tokens = tokens[1:]

		case token.CloseCurly:
			return offset, tokens, styleSheet, nil

		default:
			var (
				rule ast.Rule
				err  error
			)

			offset, tokens, rule, err = parseRule(offset, tokens)

			if err != nil {
				return offset, tokens, styleSheet, err
			}

			styleSheet.Rules = append(styleSheet.Rules, rule)
		}
	}

	return offset, tokens, styleSheet, nil
}

func parseRule(offset int, tokens []token.Token) (int, []token.Token, ast.Rule, error) {
	switch t := tokens[0].(type) {
	case token.AtKeyword:
		switch t.Value {
		case "import":
			return parseImportRule(offset+1, tokens[1:])
		case "media":
			return parseMediaRule(offset+1, tokens[1:])
		}

	default:
		return parseStyleRule(offset, tokens)
	}

	return offset, tokens, nil, SyntaxError{
		Offset:  offset,
		Message: "unexpected token",
	}
}

func parseImportRule(offset int, tokens []token.Token) (int, []token.Token, ast.ImportRule, error) {
	rule := ast.ImportRule{}

	offset, tokens = skipWhitespace(offset, tokens)

	if len(tokens) == 0 {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected end of file, expected url",
		}
	}

	switch t := tokens[0].(type) {
	case token.Url:
		parsed, err := url.Parse(t.Value)

		if err != nil {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: err.Error(),
			}
		}

		rule.URL = parsed
		tokens = tokens[1:]
		offset++

	case token.String:
		parsed, err := url.Parse(t.Value)

		if err != nil {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: err.Error(),
			}
		}

		rule.URL = parsed
		tokens = tokens[1:]
		offset++

	case token.Function:
		if t.Value != "url" {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected function, expected url()",
			}
		}

		offset, tokens = skipWhitespace(offset+1, tokens[1:])

		if len(tokens) == 0 {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected end of file",
			}
		}

		if t, ok := tokens[0].(token.String); ok {
			parsed, err := url.Parse(t.Value)

			if err != nil {
				return offset, tokens, rule, SyntaxError{
					Offset:  offset,
					Message: err.Error(),
				}
			}

			rule.URL = parsed
			tokens = tokens[1:]
			offset++
		} else {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected string",
			}
		}

		if len(tokens) == 0 {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected end of file",
			}
		}

		if _, ok := tokens[0].(token.CloseParen); ok {
			tokens = tokens[1:]
			offset++
		} else {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected closing paren",
			}
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.Semicolon); ok {
			tokens = tokens[1:]
			offset++
		} else {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected semicolon",
			}
		}
	}

	return offset, tokens, rule, nil
}

func parseMediaRule(offset int, tokens []token.Token) (int, []token.Token, ast.MediaRule, error) {
	rule := ast.MediaRule{}

	if len(tokens) == 0 {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected end of file, expected media query",
		}
	}

	offset, tokens, conditions, err := parseMediaQueryList(offset, tokens)

	if err != nil {
		return offset, tokens, rule, err
	}

	rule.Conditions = conditions

	offset, tokens = skipWhitespace(offset, tokens)

	if len(tokens) == 0 {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected end of file, expected block",
		}
	}

	if _, ok := tokens[0].(token.OpenCurly); ok {
		tokens = tokens[1:]
		offset++
	} else {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected opening curly",
		}
	}

	offset, tokens, styleSheet, err := parseStyleSheet(offset, tokens)

	if err != nil {
		return offset, tokens, rule, err
	}

	rule.StyleSheet = styleSheet

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.CloseCurly); ok {
			tokens = tokens[1:]
			offset++
		} else {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected closing curly",
			}
		}
	}

	return offset, tokens, rule, nil
}

func parseStyleRule(offset int, tokens []token.Token) (int, []token.Token, ast.StyleRule, error) {
	rule := ast.StyleRule{}

	if len(tokens) == 0 {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected end of file, expected selector",
		}
	}

	offset, tokens, selectors, err := parseSelectorList(offset, tokens)

	if err != nil {
		return offset, tokens, rule, err
	}

	rule.Selectors = selectors

	offset, tokens = skipWhitespace(offset, tokens)

	if len(tokens) == 0 {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected end of file, expected block",
		}
	}

	if _, ok := tokens[0].(token.OpenCurly); ok {
		tokens = tokens[1:]
		offset++
	} else {
		return offset, tokens, rule, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected opening curly",
		}
	}

	offset, tokens, declarations, err := parseDeclarationList(offset, tokens)

	if err != nil {
		return offset, tokens, rule, err
	}

	rule.Declarations = declarations

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.CloseCurly); ok {
			tokens = tokens[1:]
			offset++
		} else {
			return offset, tokens, rule, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected closing curly",
			}
		}
	}

	return offset, tokens, rule, nil
}

func parseDeclarationList(offset int, tokens []token.Token) (int, []token.Token, []ast.Declaration, error) {
	declarations := []ast.Declaration{}

	for len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Whitespace, token.Semicolon:
			tokens = tokens[1:]
			offset++

		case token.CloseCurly:
			return offset, tokens, declarations, nil

		case token.Ident:
			var (
				declaration ast.Declaration
				err         error
			)

			offset, tokens, declaration, err = parseDeclaration(offset+1, tokens[1:], t.Value)

			if err != nil {
				return offset, tokens, declarations, err
			}

			declarations = append(declarations, declaration)

		default:
			return offset, tokens, declarations, SyntaxError{
				Offset:  offset,
				Message: "unexpected token",
			}
		}
	}

	return offset, tokens, declarations, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseDeclaration(offset int, tokens []token.Token, name string) (int, []token.Token, ast.Declaration, error) {
	declaration := ast.Declaration{Name: name}

	offset, tokens = skipWhitespace(offset, tokens)

	if len(tokens) > 0 {
		if _, ok := tokens[0].(token.Colon); ok {
			tokens = tokens[1:]
			offset++
		} else {
			return offset, tokens, declaration, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected colon",
			}
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	for len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.CloseCurly:
			return offset, tokens, declaration, nil

		case token.Semicolon:
			return offset + 1, tokens[1:], declaration, nil

		default:
			declaration.Value = append(declaration.Value, t)
			tokens = tokens[1:]
			offset++
		}
	}

	return offset, tokens, declaration, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseSelectorList(offset int, tokens []token.Token) (int, []token.Token, []ast.Selector, error) {
	selectors := []ast.Selector{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Whitespace, token.Comma:
			tokens = tokens[1:]
			offset++

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

	return offset, tokens, selectors, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseSelector(offset int, tokens []token.Token) (int, []token.Token, ast.Selector, error) {
	var (
		left  ast.Selector
		right ast.Selector
		err   error
	)

	for len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Comma, token.OpenCurly:
			return offset, tokens, left, nil

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
				offset, tokens, right = offset+1, tokens[1:], ast.TypeSelector{Name: "*"}

				left = combineSelectors(left, right)

			case '>', '~', '+':
				offset, tokens, left, err = parseRelativeSelector(offset, tokens, left)

				if err != nil {
					return offset, tokens, left, err
				}

			default:
				return offset, tokens, left, SyntaxError{
					Offset:  offset,
					Message: "unexpected token",
				}
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
				tokens = tokens[1:]
				offset++
			}

		default:
			return offset, tokens, left, SyntaxError{
				Offset:  offset,
				Message: "unexpected token",
			}
		}
	}

	return offset, tokens, left, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func startsSelector(t token.Token) bool {
	switch t := t.(type) {
	case token.Ident:
	case token.Colon:
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

	return ast.CompoundSelector{Left: left, Right: right}
}

func parseIdSelector(offset int, tokens []token.Token) (int, []token.Token, ast.IdSelector, error) {
	selector := ast.IdSelector{}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Ident:
			selector.Name = t.Value
			return offset + 1, tokens[1:], selector, nil

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected id",
			}
		}
	}

	return offset, tokens, selector, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseClassSelector(offset int, tokens []token.Token) (int, []token.Token, ast.ClassSelector, error) {
	selector := ast.ClassSelector{}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Ident:
			selector.Name = t.Value
			return offset + 1, tokens[1:], selector, nil

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected class",
			}
		}
	}

	return offset, tokens, selector, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseAttributeSelector(offset int, tokens []token.Token) (int, []token.Token, ast.AttributeSelector, error) {
	selector := ast.AttributeSelector{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		selector.Name = t.Value
		offset, tokens = advance(offset, tokens, 1)

	default:
		return offset, tokens, selector, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected attribute name",
		}
	}

	switch t := peek(tokens, 1).(type) {
	case token.CloseSquare:
		return offset + 1, tokens[1:], selector, nil

	case token.Delim:
		switch t.Value {
		case '=':
			selector.Matcher = ast.MatcherEqual
			offset, tokens = advance(offset, tokens, 1)

		case '~', '|', '^', '$', '*':
			offset, tokens = advance(offset, tokens, 1)

			switch u := peek(tokens, 1).(type) {
			case token.Delim:
				if u.Value == '=' {
					switch t.Value {
					case '~':
						selector.Matcher = ast.MatcherIncludes
					case '|':
						selector.Matcher = ast.MatcherDashMatch
					case '^':
						selector.Matcher = ast.MatcherPrefix
					case '$':
						selector.Matcher = ast.MatcherSuffix
					case '*':
						selector.Matcher = ast.MatcherSubstring
					}

					offset, tokens = advance(offset, tokens, 1)
				}
			}
		}
	}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.String:
			selector.Value = t.Value
			tokens = tokens[1:]
			offset++

		case token.Ident:
			selector.Value = t.Value
			tokens = tokens[1:]
			offset++

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token",
			}
		}
	}

	if len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.CloseSquare:
			return offset + 1, tokens[1:], selector, nil

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected end of attribute selector",
			}
		}
	}

	return offset, tokens, selector, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseTypeSelector(offset int, tokens []token.Token) (int, []token.Token, ast.TypeSelector, error) {
	selector := ast.TypeSelector{}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Ident:
			selector.Name = t.Value
			return offset + 1, tokens[1:], selector, nil

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected type",
			}
		}
	}

	return offset, tokens, selector, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parsePseudoSelector(offset int, tokens []token.Token) (int, []token.Token, ast.PseudoSelector, error) {
	selector := ast.PseudoSelector{Name: ":"}

	if len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Colon:
			selector.Name += ":"
			tokens = tokens[1:]
			offset++
		}
	}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Ident:
			selector.Name += t.Value
			return offset + 1, tokens[1:], selector, nil

		case token.Function:
			selector.Name += t.Value
			tokens = tokens[1:]
			offset++

			for len(tokens) > 0 {
				switch t := tokens[0].(type) {
				case token.CloseParen:
					return offset + 1, tokens[1:], selector, nil
				default:
					selector.Value = append(selector.Value, t)
					tokens = tokens[1:]
					offset++
				}
			}

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected ident",
			}
		}
	}

	return offset, tokens, selector, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseRelativeSelector(offset int, tokens []token.Token, left ast.Selector) (int, []token.Token, ast.RelativeSelector, error) {
	selector := ast.RelativeSelector{Left: left}

	if len(tokens) > 0 {
		switch t := tokens[0].(type) {
		case token.Whitespace:
			selector.Combinator = ast.CombinatorDescendant

		case token.Delim:
			switch t.Value {
			case '>':
				selector.Combinator = ast.CombinatorDirectDescendant
			case '~':
				selector.Combinator = ast.CombinatorSibling
			case '+':
				selector.Combinator = ast.CombinatorDirectSibling
			default:
				return offset, tokens, selector, SyntaxError{
					Offset:  offset,
					Message: "unexpected token, expected selector combinator",
				}
			}

		default:
			return offset, tokens, selector, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected selector combinator",
			}
		}

		offset, tokens = skipWhitespace(offset+1, tokens[1:])
	}

	var (
		right ast.Selector
		err   error
	)

	offset, tokens, right, err = parseSelector(offset, tokens)

	if err != nil {
		return offset, tokens, selector, err
	}

	selector.Right = right

	return offset, tokens, selector, nil
}

func parseMediaQueryList(offset int, tokens []token.Token) (int, []token.Token, []ast.MediaQuery, error) {
	mediaQueries := []ast.MediaQuery{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Whitespace, token.Comma:
			offset, tokens = advance(offset, tokens, 1)

		case token.OpenCurly:
			return offset, tokens, mediaQueries, nil

		default:
			var (
				mediaQuery ast.MediaQuery
				err        error
			)

			offset, tokens, mediaQuery, err = parseMediaQuery(offset, tokens)

			if err != nil {
				return offset, tokens, mediaQueries, err
			}

			mediaQueries = append(mediaQueries, mediaQuery)
		}
	}

	return offset, tokens, mediaQueries, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func parseMediaQuery(offset int, tokens []token.Token) (int, []token.Token, ast.MediaQuery, error) {
	mediaQuery := ast.MediaQuery{}

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		switch t.Value {
		case "not":
			mediaQuery.Qualifier = ast.QualifierNot
			offset, tokens = advance(offset, tokens, 1)

		case "only":
			mediaQuery.Qualifier = ast.QualifierOnly
			offset, tokens = advance(offset, tokens, 1)
		}

	case token.OpenParen:
		var (
			condition ast.MediaCondition
			err       error
		)

		offset, tokens, condition, err = parseMediaCondition(offset, tokens)

		if err != nil {
			return offset, tokens, mediaQuery, err
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
			return offset, tokens, mediaQuery, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected media type",
			}

		default:
			mediaQuery.Type = t.Value
			offset, tokens = advance(offset, tokens, 1)
		}

	case nil:
		return offset, tokens, mediaQuery, SyntaxError{
			Offset:  offset,
			Message: "unexpected end of file",
		}

	default:
		return offset, tokens, mediaQuery, SyntaxError{
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
				err error
			)

			offset, tokens, condition, err = parseMediaExpression(
				skipWhitespace(offset+1, tokens[1:]),
			)

			if err != nil {
				return offset, tokens, nil, err
			}

			condition = ast.MediaNegation{Condition: condition}

			return offset, tokens, condition, nil
		}

		return offset, tokens, nil, SyntaxError{
			Offset: offset,
			Message: "unexpected token",
		}

	default:
		return parseMediaExpression(offset, tokens)
	}
}

func parseMediaExpression(offset int, tokens []token.Token) (int, []token.Token, ast.MediaCondition, error) {
	if offset, tokens, feature, err := parseMediaFeature(offset, tokens); err == nil {
		return offset, tokens, feature, nil
	}

	return offset, tokens, nil, SyntaxError{
		Offset: offset,
		Message: "Shit, Sherlock",
	}
}

func parseMediaFeature(offset int, tokens []token.Token) (int, []token.Token, ast.MediaFeature, error) {
	feature := ast.MediaFeature{}

	switch peek(tokens, 1).(type) {
	case token.OpenParen:
		offset, tokens = advance(offset, tokens, 1)
	default:
		return offset, tokens, feature, SyntaxError{
			Offset: offset,
			Message: "unexpected token, expected opening paren",
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Ident:
		feature.Name = t.Value
		offset, tokens = advance(offset, tokens, 1)
	default:
		return offset, tokens, feature, SyntaxError{
			Offset: offset,
			Message: "unexpected token, expected feature name",
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.Colon:
		offset, tokens = advance(offset, tokens, 1)
	default:
		return offset, tokens, feature, SyntaxError{
			Offset: offset,
			Message: "unexpected token, expected colon",
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch t := peek(tokens, 1).(type) {
	case token.Number:
		feature.Value = ast.MediaValuePlain{Value: t}
		offset, tokens = advance(offset, tokens, 1)

	case token.Dimension:
		feature.Value = ast.MediaValuePlain{Value: t}
		offset, tokens = advance(offset, tokens, 1)

	case token.Ident:
		feature.Value = ast.MediaValuePlain{Value: t}
		offset, tokens = advance(offset, tokens, 1)

	default:
		return offset, tokens, feature, SyntaxError{
			Offset: offset,
			Message: "unexpected token, expected colon",
		}
	}

	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.CloseParen:
		offset, tokens = advance(offset, tokens, 1)
	default:
		return offset, tokens, feature, SyntaxError{
			Offset: offset,
			Message: "unexpected token, expected opening paren",
		}
	}

	return offset, tokens, feature, nil
}

func peek(tokens []token.Token, n int) token.Token {
	if len(tokens) < n {
		return nil
	}

	return tokens[n - 1]
}

func advance(offset int, tokens []token.Token, n int) (int, []token.Token) {
	return offset + n, tokens[n:]
}

func skipWhitespace(offset int, tokens []token.Token) (int, []token.Token) {
	if _, ok := peek(tokens, 1).(token.Whitespace); ok {
		return advance(offset, tokens, 1)
	}

	return offset, tokens
}
