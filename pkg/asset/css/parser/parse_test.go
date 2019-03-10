package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input      string
		styleSheet ast.StyleSheet
	}{
		{
			".foo{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.ClassSelector{Name: "foo"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"#foo{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.IdSelector{Name: "foo"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"#foo,.bar{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.IdSelector{Name: "foo"},
							ast.ClassSelector{Name: "bar"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"#foo.bar{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.CompoundSelector{
								Left:  ast.IdSelector{Name: "foo"},
								Right: ast.ClassSelector{Name: "bar"},
							},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"#foo.bar.baz{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.CompoundSelector{
								Left: ast.CompoundSelector{
									Left:  ast.IdSelector{Name: "foo"},
									Right: ast.ClassSelector{Name: "bar"},
								},
								Right: ast.ClassSelector{Name: "baz"},
							},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"#foo.bar>.baz{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.RelativeSelector{
								Combinator: ast.DirectDescendant,
								Left: ast.CompoundSelector{
									Left:  ast.IdSelector{Name: "foo"},
									Right: ast.ClassSelector{Name: "bar"},
								},
								Right: ast.ClassSelector{Name: "baz"},
							},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"#foo.bar .baz{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.RelativeSelector{
								Combinator: ast.Descendant,
								Left: ast.CompoundSelector{
									Left:  ast.IdSelector{Name: "foo"},
									Right: ast.ClassSelector{Name: "bar"},
								},
								Right: ast.ClassSelector{Name: "baz"},
							},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			`@import "foo"`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{Url: "foo"},
				},
			},
		},
		{
			`@import url(foo)`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{Url: "foo"},
				},
			},
		},
		{
			`@import url("foo")`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{Url: "foo"},
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
