package parser

import (
	"testing"
	"net/url"

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
								Combinator: ast.CombinatorDirectDescendant,
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
								Combinator: ast.CombinatorDescendant,
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
			":foo{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.PseudoSelector{Name: ":foo"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"::foo{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.PseudoSelector{Name: "::foo"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"*{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.TypeSelector{Name: "*"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"[foo]{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.AttributeSelector{Name: "foo"},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			"[foo=bar]{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.AttributeSelector{
								Name: "foo",
								Matcher: ast.MatcherEqual,
								Value: "bar",
							},
						},
						Declarations: []ast.Declaration{},
					},
				},
			},
		},
		{
			`[foo="bar"]{}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.AttributeSelector{
								Name: "foo",
								Matcher: ast.MatcherEqual,
								Value: "bar",
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
					ast.ImportRule{URL: &url.URL{Path: "foo"}},
				},
			},
		},
		{
			`@import url(foo)`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{URL: &url.URL{Path: "foo"}},
				},
			},
		},
		{
			`@import url("foo")`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.ImportRule{URL: &url.URL{Path: "foo"}},
				},
			},
		},
		{
			`@media screen {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.MediaRule{
						Conditions: []ast.MediaQuery{
							ast.MediaQuery{Type: "screen"},
						},
					},
				},
			},
		},
		{
			`@media only screen {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.MediaRule{
						Conditions: []ast.MediaQuery{
							ast.MediaQuery{Type: "screen", Qualifier: ast.QualifierOnly},
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
