package parser

import (
	"net/url"
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
			"#foo{}",
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.StyleRule{
						Selectors: []ast.Selector{
							ast.IdSelector{Name: "foo"},
						},
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
								Combinator: '>',
								Left: ast.CompoundSelector{
									Left:  ast.IdSelector{Name: "foo"},
									Right: ast.ClassSelector{Name: "bar"},
								},
								Right: ast.ClassSelector{Name: "baz"},
							},
						},
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
								Combinator: ' ',
								Left: ast.CompoundSelector{
									Left:  ast.IdSelector{Name: "foo"},
									Right: ast.ClassSelector{Name: "bar"},
								},
								Right: ast.ClassSelector{Name: "baz"},
							},
						},
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
								Name:    "foo",
								Matcher: "=",
								Value:   "bar",
							},
						},
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
								Name:    "foo",
								Matcher: "=",
								Value:   "bar",
							},
						},
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
							ast.MediaQuery{Type: "screen", Qualifier: "only"},
						},
					},
				},
			},
		},
		{
			`@media screen, print {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.MediaRule{
						Conditions: []ast.MediaQuery{
							ast.MediaQuery{Type: "screen"},
							ast.MediaQuery{Type: "print"},
						},
					},
				},
			},
		},
		{
			`@media (foo: bar) {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.MediaRule{
						Conditions: []ast.MediaQuery{
							ast.MediaQuery{
								Condition: ast.MediaFeature{
									Name: "foo",
									Value: ast.MediaValuePlain{Value: token.Ident{
										Offset: 13,
										Value:  "bar",
									}},
								},
							},
						},
					},
				},
			},
		},
		{
			`@media ((foo: bar) and (baz: qux)) {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.MediaRule{
						Conditions: []ast.MediaQuery{
							ast.MediaQuery{
								Condition: ast.MediaOperation{
									Operator: "and",
									Left: ast.MediaFeature{
										Name: "foo",
										Value: ast.MediaValuePlain{Value: token.Ident{
											Offset: 14,
											Value:  "bar",
										}},
									},
									Right: ast.MediaFeature{
										Name: "baz",
										Value: ast.MediaValuePlain{Value: token.Ident{
											Offset: 29,
											Value:  "qux",
										}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			`@media ((foo: bar) and ((baz: qux) or (fez: fud))) {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.MediaRule{
						Conditions: []ast.MediaQuery{
							ast.MediaQuery{
								Condition: ast.MediaOperation{
									Operator: "and",
									Left: ast.MediaFeature{
										Name: "foo",
										Value: ast.MediaValuePlain{Value: token.Ident{
											Offset: 14,
											Value:  "bar",
										}},
									},
									Right: ast.MediaOperation{
										Operator: "or",
										Left: ast.MediaFeature{
											Name: "baz",
											Value: ast.MediaValuePlain{Value: token.Ident{
												Offset: 30,
												Value:  "qux",
											}},
										},
										Right: ast.MediaFeature{
											Name: "fez",
											Value: ast.MediaValuePlain{Value: token.Ident{
												Offset: 44,
												Value:  "fud",
											}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			`@keyframes foo {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.KeyframesRule{Name: "foo"},
				},
			},
		},
		{
			`@keyframes "foo" {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.KeyframesRule{Name: "foo"},
				},
			},
		},
		{
			`@-webkit-keyframes foo {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.KeyframesRule{Prefix: "-webkit-", Name: "foo"},
				},
			},
		},
		{
			`@keyframes "foo" { from {} to {} }`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.KeyframesRule{
						Name: "foo",
						Blocks: []ast.KeyframeBlock{
							ast.KeyframeBlock{Selector: 0},
							ast.KeyframeBlock{Selector: 1},
						},
					},
				},
			},
		},
		{
			`@supports (foo: bar) {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.SupportsRule{
						Condition: ast.SupportsFeature{
							Declaration: ast.Declaration{
								Name: "foo",
								Value: []token.Token{
									token.Ident{Offset: 16, Value: "bar"},
								},
							},
						},
					},
				},
			},
		},
		{
			`@page {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.PageRule{},
				},
			},
		},
		{
			`@page foo {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.PageRule{
						Selectors: []ast.PageSelector{
							ast.PageSelector{
								Type: "foo",
							},
						},
					},
				},
			},
		},
		{
			`@page foo:left {}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.PageRule{
						Selectors: []ast.PageSelector{
							ast.PageSelector{
								Type:    "foo",
								Classes: []string{":left"},
							},
						},
					},
				},
			},
		},
		{
			`@page foo:left {color:red}`,
			ast.StyleSheet{
				Rules: []ast.Rule{
					ast.PageRule{
						Selectors: []ast.PageSelector{
							ast.PageSelector{
								Type:    "foo",
								Classes: []string{":left"},
							},
						},
						Components: []ast.PageComponent{
							ast.PageDeclaration{
								Declaration: ast.Declaration{
									Name: "color",
									Value: []token.Token{
										token.Ident{Offset: 22, Value: "red"},
									},
								},
							},
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
