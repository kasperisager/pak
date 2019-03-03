package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input      string
		styleSheet ast.StyleSheet
	}{
		{
			".foo",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.QualifiedRule{
						Prelude: []token.Token{
							token.Delim{Offset: 0, Value: '.'},
							token.Ident{Offset: 1, Value: "foo"},
						},
					},
				},
			},
		},
		{
			"#foo",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.QualifiedRule{
						Prelude: []token.Token{
							token.Hash{Offset: 0, Value: "foo", Id: true},
						},
					},
				},
			},
		},
		{
			"@foo bar",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.AtRule{
						Name: "foo",
						Prelude: []token.Token{
							token.Whitespace{Offset: 4},
							token.Ident{Offset: 5, Value: "bar"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		runes := []rune(test.input)

		tokens, err := scanner.Scan(runes)
		assert.Nil(t, err, test.input)

		ast, err := Parse(tokens)
		assert.Nil(t, err, test.input)

		assert.Equal(t, test.styleSheet, ast, test.input)
	}
}
