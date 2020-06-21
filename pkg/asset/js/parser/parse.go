package parser

import (
	"github.com/kasperisager/pak/pkg/asset/js/ast"
	"github.com/kasperisager/pak/pkg/asset/js/scanner"
	"github.com/kasperisager/pak/pkg/asset/js/token"
)

type SyntaxError struct {
	Offset  int
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

func Parse(runes []rune) (*ast.Program, error) {
	parameters := parameters{
		In:     true,
		Yield:  false,
		Await:  false,
		Tagged: true,
		Return: false,
	}

	parser := parser{0, runes, nil}

	program := &ast.Program{}

	for {
		var (
			statement ast.Statement
			ok        bool
			err       error
		)

		parser, statement, ok, err = parseStatement(parser, parameters)

		if err != nil {
			return program, err
		}

		if ok {
			program.Body = append(program.Body, statement)
		} else {
			break
		}
	}

	return program, nil
}

type (
	parser struct {
		offset int
		runes  []rune
		tokens []token.Token
	}

	parameters struct {
		In, Yield, Await, Tagged, Return bool
	}
)

func (p parser) scan(options ...func(*scanner.Options)) (parser, error) {
	var (
		token token.Token
		err   error
	)

	p.offset, p.runes, token, err = scanner.Scan(
		p.offset,
		p.runes,
		options...,
	)

	if err != nil {
		return p, err
	}

	p.tokens = append(p.tokens, token)

	return p, nil
}

func (p parser) scanN(n int, options ...func(*scanner.Options)) (parser, error) {
	for len(p.tokens) < n {
		var err error

		p, err = p.scan(options...)

		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func (p parser) peek(n int, options ...func(*scanner.Options)) (parser, token.Token) {
	p, err := p.scanN(n, options...)

	if err != nil {
		return p, nil
	}

	return p, p.tokens[n-1]
}

func (p parser) advance(n int, options ...func(*scanner.Options)) parser {
	p, token := p.peek(n, options...)

	if token != nil {
		p.tokens = p.tokens[n:]
	}

	return p
}

// https://www.ecma-international.org/ecma-262/#prod-Statement
func parseStatement(parser parser, parameters parameters) (parser, ast.Statement, bool, error) {
	parser, next := parser.peek(1)

	switch next := next.(type) {
	case token.Punctuator:
		switch next.Value {
		case "{":
			return parseBlockStatement(parser, parameters)
		}
	}

	return parseExpressionStatement(parser, parameters)
}

// https://www.ecma-international.org/ecma-262/#prod-BlockStatement
func parseBlockStatement(parser parser, parameters parameters) (parser, *ast.BlockStatement, bool, error) {
	blockStatement := &ast.BlockStatement{}

	parser, next := parser.peek(1)

	if next, ok := next.(token.Punctuator); ok && next.Value == "{" {
		parser = parser.advance(1)
	} else {
		return parser, nil, false, nil
	}

	for {
		var (
			statement ast.Statement
			ok        bool
			err       error
		)

		parser, statement, ok, err = parseStatement(parser, parameters)

		if err != nil {
			return parser, nil, false, err
		}

		if ok {
			blockStatement.Body = append(blockStatement.Body, statement)
		} else {
			break
		}
	}

	parser, next = parser.peek(1)

	if next, ok := next.(token.Punctuator); ok && next.Value == "}" {
		parser = parser.advance(1)
	} else {
		return parser, nil, false, SyntaxError{
			Offset:  parser.offset,
			Message: `unexpected token, expected "}"`,
		}
	}

	return parser, blockStatement, true, nil
}

// https://www.ecma-international.org/ecma-262/#prod-ExpressionStatement
func parseExpressionStatement(parser parser, parameters parameters) (parser, *ast.ExpressionStatement, bool, error) {
	parameters.In = true

	parser, expression, ok, err := parseExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		expressionStatement := &ast.ExpressionStatement{
			Expression: expression,
		}

		return parser, expressionStatement, true, nil
	}

	return parser, nil, false, nil
}

// https://www.ecma-international.org/ecma-262/#prod-Expression
func parseExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseAssignmentExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	parser, next := parser.peek(1)

	if punctuator, ok := next.(token.Punctuator); ok && punctuator.Value == "," {
		sequenceExpression := &ast.SequenceExpression{
			Expression: []ast.Expression{left},
		}

		for {
			if punctuator, ok := next.(token.Punctuator); ok && punctuator.Value == "," {
				parser = parser.advance(1)

				var (
					right ast.Expression
					ok    bool
					err   error
				)

				parser, right, ok, err = parseExpression(parser, parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if ok {
					sequenceExpression.Expression = append(sequenceExpression.Expression, right)
				} else {
					break
				}
			} else {
				break
			}
		}

		return parser, sequenceExpression, true, nil
	}

	return parser, left, ok, nil
}

// https://www.ecma-international.org/ecma-262/#prod-AssignmentExpression
func parseAssignmentExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseConditionalExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		switch next := next.(type) {
		case token.Punctuator:
			switch next.Value {
			case "=", "*=", "/=", "%=", "+=", "-=", "<<=", ">>=", ">>>=", "&=", "^=", "|=", "**=":
				left, ok := left.(ast.Pattern)

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "cannot assign to expression",
					}
				}

				parser = parser.advance(1)

				var (
					right ast.Expression
					err   error
				)

				parser, right, ok, err = parseExpression(parser, parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "unexpected assignment",
					}
				}

				assignmentExpression := &ast.AssignmentExpression{
					Operator: next.Value,
					Left:     left,
					Right:    right,
				}

				return parser, assignmentExpression, true, nil

			}
		}

		return parser, left, true, nil
	}

	if parameters.Yield {
		return parseYieldExpression(parser, parameters)
	}

	return parser, nil, false, nil
}

// https://www.ecma-international.org/ecma-262/#prod-ConditionalExpression
func parseConditionalExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, test, ok, err := parseLogicalOrExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok && next.Value == "?" {
			parser, alternate, ok, err := parseAssignmentExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected assignment expression",
				}
			}

			parser, next := parser.peek(1)

			if next, ok := next.(token.Punctuator); ok && next.Value == ":" {
				parser, consequent, ok, err := parseAssignmentExpression(parser.advance(1), parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "expected assignment expression",
					}
				}

				conditionalExpression := &ast.ConditionalExpression{
					Test:       test,
					Alternate:  alternate,
					Consequent: consequent,
				}

				return parser, conditionalExpression, true, nil
			}
		}
	}

	return parser, test, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-LogicalORExpression
func parseLogicalOrExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseLogicalAndExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok && next.Value == "||" {
			parser, right, ok, err := parseLogicalOrExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected logical OR expression",
				}
			}

			logicalExpression := &ast.LogicalExpression{
				Operator: "||",
				Left:     left,
				Right:    right,
			}

			return parser, logicalExpression, true, nil
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-LogicalANDExpression
func parseLogicalAndExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseBitwiseOrExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok && next.Value == "&&" {
			parser, right, ok, err := parseLogicalAndExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected logical AND expression",
				}
			}

			logicalExpression := &ast.LogicalExpression{
				Operator: "&&",
				Left:     left,
				Right:    right,
			}

			return parser, logicalExpression, true, nil
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-BitwiseORExpression
func parseBitwiseOrExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseBitwiseXorExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok && next.Value == "|" {
			parser, right, ok, err := parseBitwiseOrExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected bitwise OR expression",
				}
			}

			binaryExpression := &ast.BinaryExpression{
				Operator: "|",
				Left:     left,
				Right:    right,
			}

			return parser, binaryExpression, true, nil
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-BitwiseXORExpression
func parseBitwiseXorExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseBitwiseAndExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok && next.Value == "^" {
			parser, right, ok, err := parseBitwiseXorExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected bitwise XOR expression",
				}
			}

			binaryExpression := &ast.BinaryExpression{
				Operator: "^",
				Left:     left,
				Right:    right,
			}

			return parser, binaryExpression, true, nil
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-BitwiseANDExpression
func parseBitwiseAndExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseEqualityExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok && next.Value == "&" {
			parser, right, ok, err := parseBitwiseAndExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected bitwise AND expression",
				}
			}

			binaryExpression := &ast.BinaryExpression{
				Operator: "&",
				Left:     left,
				Right:    right,
			}

			return parser, binaryExpression, true, nil
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-EqualityExpression
func parseEqualityExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseRelationalExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok {
			switch next.Value {
			case "==", "!=", "===", "!==":
				parser, right, ok, err := parseEqualityExpression(parser.advance(1), parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "expected equality expression",
					}
				}

				binaryExpression := &ast.BinaryExpression{
					Operator: next.Value,
					Left:     left,
					Right:    right,
				}

				return parser, binaryExpression, true, nil
			}
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-RelationalExpression
func parseRelationalExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseShiftExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok {
			switch next.Value {
			case "<", ">", "<=", ">=":
				parser, right, ok, err := parseRelationalExpression(parser.advance(1), parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "expected relational expression",
					}
				}

				binaryExpression := &ast.BinaryExpression{
					Operator: next.Value,
					Left:     left,
					Right:    right,
				}

				return parser, binaryExpression, true, nil
			}
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-ShiftExpression
func parseShiftExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseAdditiveExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok {
			switch next.Value {
			case "<<", ">>", ">>>":
				parser, right, ok, err := parseShiftExpression(parser.advance(1), parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "expected shift expression",
					}
				}

				binaryExpression := &ast.BinaryExpression{
					Operator: next.Value,
					Left:     left,
					Right:    right,
				}

				return parser, binaryExpression, true, nil
			}
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-AdditiveExpression
func parseAdditiveExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseMultiplicativeExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok {
			switch next.Value {
			case "+", "-":
				parser, right, ok, err := parseAdditiveExpression(parser.advance(1), parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "expected additive expression",
					}
				}

				binaryExpression := &ast.BinaryExpression{
					Operator: next.Value,
					Left:     left,
					Right:    right,
				}

				return parser, binaryExpression, true, nil
			}
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-MultiplicativeExpression
func parseMultiplicativeExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseExponentiationExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok {
		parser, next := parser.peek(1)

		if next, ok := next.(token.Punctuator); ok {
			switch next.Value {
			case "*", "/", "%":
				parser, right, ok, err := parseExponentiationExpression(parser.advance(1), parameters)

				if err != nil {
					return parser, nil, false, err
				}

				if !ok {
					return parser, nil, false, SyntaxError{
						Offset:  parser.offset,
						Message: "expected exponentiation expression",
					}
				}

				binaryExpression := &ast.BinaryExpression{
					Operator: next.Value,
					Left:     left,
					Right:    right,
				}

				return parser, binaryExpression, true, nil
			}
		}
	}

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-ExponentiationExpression
func parseExponentiationExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseUnaryExpression(parser, parameters)

	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-UnaryExpression
func parseUnaryExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, next := parser.peek(1)

	var operator string

	if next, ok := next.(token.Punctuator); ok {
		operator = next.Value

		switch operator {
		case "+", "-":
			parser, next := parser.peek(2)

			if next, ok := next.(token.Punctuator); ok && next.Value == operator {
				return parseUpdateExpression(parser, parameters)
			}

		case "~", "!":

		default:
			return parseUpdateExpression(parser, parameters)
		}

		parser = parser.advance(1)
	}

	if next, ok := next.(token.Keyword); ok {
		operator = next.Value

		switch operator {
		case "delete", "void", "typeof":

		default:
			return parseUpdateExpression(parser, parameters)
		}

		parser = parser.advance(1)
	}

	parser, argument, ok, err := parseUpdateExpression(parser, parameters)

	if err != nil {
		return parser, nil, false, err
	}

	if ok && operator != "" {
		unaryExpression := &ast.UnaryExpression{
			Operator: operator,
			Prefix:   true,
			Argument: argument,
		}

		return parser, unaryExpression, true, nil
	}

	return parser, argument, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-UpdateExpression
func parseUpdateExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, next := parser.peek(1)

	if next, ok := next.(token.Punctuator); ok {
		operator := next.Value

		switch operator {
		case "++", "--":
			parser, argument, ok, err := parseUnaryExpression(parser.advance(1), parameters)

			if err != nil {
				return parser, nil, false, err
			}

			if !ok {
				return parser, nil, false, SyntaxError{
					Offset:  parser.offset,
					Message: "expected unary expression",
				}
			}

			updateExpression := &ast.UpdateExpression{
				Operator: operator,
				Prefix:   true,
				Argument: argument,
			}

			return parser, updateExpression, true, nil
		}
	}

	parser, argument, ok, err := parseLeftHandSideExpression(parser, parameters)

	parser, next = parser.peek(1)

	if next, ok := next.(token.Punctuator); ok {
		operator := next.Value

		switch operator {
		case "++", "--":
			updateExpression := &ast.UpdateExpression{
				Operator: operator,
				Prefix:   false,
				Argument: argument,
			}

			return parser.advance(1), updateExpression, true, nil
		}
	}

	return parser, argument, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-LeftHandSideExpression
func parseLeftHandSideExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseNewExpression(parser, parameters)
	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-NewExpression
func parseNewExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parseMemberExpression(parser, parameters)
	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-MemberExpression
func parseMemberExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, left, ok, err := parsePrimaryExpression(parser, parameters)
	return parser, left, ok, err
}

// https://www.ecma-international.org/ecma-262/#prod-CallExpression
func parseCallExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	return parser, nil, false, nil
}

// https://www.ecma-international.org/ecma-262/#prod-YieldExpression
func parseYieldExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	return parser, nil, false, nil
}

// https://www.ecma-international.org/ecma-262/#prod-PrimaryExpression
func parsePrimaryExpression(parser parser, parameters parameters) (parser, ast.Expression, bool, error) {
	parser, next := parser.peek(1)

	switch next := next.(type) {
	case token.String:
		return parser.advance(1), &ast.StringLiteral{Value: next.Value}, true, nil

	case token.Number:
		return parser.advance(1), &ast.NumberLiteral{Value: next.Value}, true, nil

	case token.Boolean:
		return parser.advance(1), &ast.BooleanLiteral{Value: next.Value}, true, nil

	case token.Null:
		return parser.advance(1), &ast.NullLiteral{}, true, nil

	case token.Identifier:
		return parseIdentifierReference(parser, parameters)
	}

	return parser, nil, false, nil
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierReference
func parseIdentifierReference(parser parser, parameters parameters) (parser, *ast.Identifier, bool, error) {
	parser, next := parser.peek(1)

	switch next := next.(type) {
	case token.Identifier:
		switch next.Value {
		case "yield":
			if !parameters.Yield {
				return parser.advance(1), &ast.Identifier{Name: "yield"}, true, nil
			}

		case "await":
			if !parameters.Await {
				return parser.advance(1), &ast.Identifier{Name: "await"}, true, nil
			}

		default:
			return parser.advance(1), &ast.Identifier{Name: next.Value}, true, nil
		}
	}

	return parser, nil, false, nil
}
