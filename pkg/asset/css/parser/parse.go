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
	switch token := tokens[0].(type) {
	case token.AtKeyword:
		return parseAtRule(tokens[1:], token.Value)

	default:
		return parseQualifiedRule(tokens)
	}
}

func parseAtRule(tokens []token.Token, name string) ([]token.Token, ast.AtRule, error) {
	rule := ast.AtRule{Name: name}

	for len(tokens) > 0 {
		switch tokens[0].(type) {
		case token.Semicolon:
			rule.Prelude = token.TrimWhitespace(rule.Prelude)

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
			rule.Prelude = token.TrimWhitespace(rule.Prelude)

			return tokens, rule, nil

		default:
			rule.Prelude = append(rule.Prelude, tokens[0])
			tokens = tokens[1:]
		}
	}

	rule.Prelude = token.TrimWhitespace(rule.Prelude)

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
			rule.Prelude = token.TrimWhitespace(rule.Prelude)

			return tokens, rule, nil

		default:
			rule.Prelude = append(rule.Prelude, tokens[0])
			tokens = tokens[1:]
		}
	}

	rule.Prelude = token.TrimWhitespace(rule.Prelude)

	return tokens, rule, nil
}

func parseBlock(tokens []token.Token) ([]token.Token, ast.Block, error) {
	block := ast.Block{}
	end := 0

	for len(tokens) > end {
		switch tokens[end].(type) {
		case token.CloseCurly:
			block.Value = token.TrimWhitespace(tokens[0:end])

			return tokens[end+1:], block, nil

		default:
			end++
		}
	}

	return tokens[end:], block, nil
}
