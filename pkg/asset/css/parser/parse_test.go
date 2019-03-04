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
							token.Ident{Offset: 5, Value: "bar"},
						},
					},
				},
			},
		},
		{
			`@import "foo"`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{
						AtRule: ast.AtRule{
							Name: "import",
							Prelude: []token.Token{
								token.String{Offset: 8, Value: "foo"},
							},
						},
						Url: "foo",
					},
				},
			},
		},
		{
			`@import url(foo)`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{
						AtRule: ast.AtRule{
							Name: "import",
							Prelude: []token.Token{
								token.Url{Offset: 8, Value: "foo"},
							},
						},
						Url: "foo",
					},
				},
			},
		},
		{
			`@import url("foo")`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{
						AtRule: ast.AtRule{
							Name: "import",
							Prelude: []token.Token{
								token.Function{Offset: 8, Value: "url"},
								token.String{Offset: 12, Value: "foo"},
								token.CloseParen{Offset: 17},
							},
						},
						Url: "foo",
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
