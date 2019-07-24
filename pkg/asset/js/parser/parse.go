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

type Scanner struct {
	offset int
	runes  []rune
	tokens []token.Token
}

func (s *Scanner) Offset() int {
	return s.offset
}

func (s *Scanner) Peek(n int, options ...func(*scanner.Options)) token.Token {
	for len(s.tokens) < n {
		var ok bool

		s.offset, s.runes, s.tokens, ok = scanner.ScanInto(s.offset, s.runes, s.tokens, options...)

		if !ok {
			return nil
		}
	}

	return s.tokens[n-1]
}

func (s *Scanner) Advance(n int, options ...func(*scanner.Options)) token.Token {
	token := s.Peek(n, options...)

	if token == nil {
		return nil
	}

	s.tokens = s.tokens[n:]

	return token
}

type parameters struct {
	In, Yield, Await, Tagged, Return bool
}

func Parse(runes []rune) (*ast.Program, error) {
	parameters := parameters{
		In:     true,
		Yield:  false,
		Await:  false,
		Tagged: true,
		Return: false,
	}

	scanner := &Scanner{0, runes, nil}

	var program *ast.Program

	for {
		statement, ok, err := parseStatement(scanner, parameters)

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

// https://www.ecma-international.org/ecma-262/#prod-Statement
func parseStatement(scanner *Scanner, parameters parameters) (ast.Statement, bool, error) {
	block, ok, err := parseBlockStatement(scanner, parameters)

	if err != nil {
		return nil, false, err
	}

	if ok {
		return block, true, nil
	}

	return nil, false, nil
}

// https://www.ecma-international.org/ecma-262/#prod-BlockStatement
func parseBlockStatement(scanner *Scanner, parameters parameters) (*ast.BlockStatement, bool, error) {
	block := &ast.BlockStatement{}

	next := scanner.Peek(1)

	if next, ok := next.(token.Punctuator); ok && next.Value == "{" {
		scanner.Advance(1)
	} else {
		return block, false, nil
	}

	for {
		statement, ok, err := parseStatement(scanner, parameters)

		if err != nil {
			return block, false, err
		}

		if ok {
			block.Body = append(block.Body, statement)
		} else {
			break
		}
	}

	next = scanner.Peek(1)

	if next, ok := next.(token.Punctuator); ok && next.Value == "}" {
		scanner.Advance(1)
	} else {
		return block, false, SyntaxError{
			Offset:  scanner.Offset(),
			Message: "unexpected token",
		}
	}

	return block, true, nil
}

// https://www.ecma-international.org/ecma-262/#prod-PrimaryExpression
// func parsePrimaryExpression(scanner *Scanner, parameters parameters) (ast.Expression, bool, error) {
// }
