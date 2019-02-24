package parser

import (
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

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
	switch tokens[0].(type) {
	case token.AtKeyword:
		return parseAtRule(tokens)

	default:
		return parseQualifiedRule(tokens)
	}
}

func parseAtRule(tokens []token.Token) ([]token.Token, ast.AtRule, error) {
	rule := ast.AtRule{Name: tokens[0].(token.AtKeyword).Value}

	tokens = tokens[1:]

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Semicolon:
			return tokens[1:], rule, nil

		case token.OpenCurly:
			var (
				block ast.Block
				err   error
			)

			tokens, block, err = parseBlock(tokens[1:])

			if err != nil {
				return tokens, rule, err
			}

			rule.Value = &block

			return tokens, rule, nil

		default:
			var (
				preserved ast.Preserved
				err       error
			)

			tokens, preserved, err = parsePreserved(tokens)

			if err != nil {
				return tokens, rule, err
			}

			rule.Prelude = append(rule.Prelude, preserved)
		}
	}

	return tokens, rule, nil
}

func parseQualifiedRule(tokens []token.Token) ([]token.Token, ast.QualifiedRule, error) {
	rule := ast.QualifiedRule{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.OpenCurly:
			var (
				block ast.Block
				err   error
			)

			tokens, block, err = parseBlock(tokens[1:])

			if err != nil {
				return tokens, rule, err
			}

			rule.Value = block

			return tokens, rule, nil

		default:
			var (
				preserved ast.Preserved
				err       error
			)

			tokens, preserved, err = parsePreserved(tokens)

			if err != nil {
				return tokens, rule, err
			}

			rule.Prelude = append(rule.Prelude, preserved)
		}
	}

	return tokens, rule, nil
}

func parseBlock(tokens []token.Token) ([]token.Token, ast.Block, error) {
	block := ast.Block{}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.CloseCurly:
			return tokens[1:], block, nil

		default:
			var (
				preserved ast.Preserved
				err       error
			)

			tokens, preserved, err = parsePreserved(tokens)

			if err != nil {
				return tokens, block, err
			}

			block.Value = append(block.Value, preserved)
		}
	}

	return tokens, block, nil
}

func parsePreserved(tokens []token.Token) ([]token.Token, ast.Preserved, error) {
	switch token := tokens[0].(type) {
	case token.Ident:
		return tokens[1:], ast.Ident{Ident: token}, nil

	case token.AtKeyword:
		return tokens[1:], ast.AtKeyword{AtKeyword: token}, nil

	case token.Hash:
		return tokens[1:], ast.Hash{Hash: token}, nil

	case token.String:
		return tokens[1:], ast.String{String: token}, nil

	case token.Url:
		return tokens[1:], ast.Url{Url: token}, nil

	case token.Delim:
		return tokens[1:], ast.Delim{Delim: token}, nil

	case token.Number:
		return tokens[1:], ast.Number{Number: token}, nil

	case token.Percentage:
		return tokens[1:], ast.Percentage{Percentage: token}, nil

	case token.Dimension:
		return tokens[1:], ast.Dimension{Dimension: token}, nil

	case token.Whitespace:
		return tokens[1:], ast.Whitespace{Whitespace: token}, nil

	case token.Colon:
		return tokens[1:], ast.Colon{Colon: token}, nil

	case token.Semicolon:
		return tokens[1:], ast.Semicolon{Semicolon: token}, nil

	case token.Comma:
		return tokens[1:], ast.Comma{Comma: token}, nil

	case token.CloseSquare:
		return tokens[1:], ast.CloseSquare{CloseSquare: token}, nil

	case token.CloseParen:
		return tokens[1:], ast.CloseParen{CloseParen: token}, nil

	case token.CloseCurly:
		return tokens[1:], ast.CloseCurly{CloseCurly: token}, nil
	}

	return tokens, nil, nil
}
