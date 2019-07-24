package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/html/ast"
	"github.com/kasperisager/pak/pkg/asset/html/scanner"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input string
		root  ast.Element
	}{
		{
			`
			<!doctype html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html></html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html class="foo"></html>
			`,
			ast.Element{
				Name: "html",
				Attributes: []ast.Attribute{
					ast.Attribute{Name: "class", Value: "foo"},
				},
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<head></head>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<head class="foo"></head>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{
						Name: "head",
						Attributes: []ast.Attribute{
							ast.Attribute{Name: "class", Value: "foo"},
						},
					},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<head>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			</head>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<head>
					<link rel="stylesheet" href="foo.css">
				</head>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{
						Name: "head",
						Children: []ast.Node{
							ast.Element{
								Name: "link",
								Attributes: []ast.Attribute{
									ast.Attribute{Name: "rel", Value: "stylesheet"},
									ast.Attribute{Name: "href", Value: "foo.css"},
								},
							},
						},
					},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<head>
					<title>Foo</title>
				</head>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{
						Name: "head",
						Children: []ast.Node{
							ast.Element{
								Name: "title",
								Children: []ast.Node{
									ast.Text{Data: "Foo"},
								},
							},
						},
					},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<body></body>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<body class="foo"></body>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{
						Name: "body",
						Attributes: []ast.Attribute{
							ast.Attribute{Name: "class", Value: "foo"},
						},
					},
				},
			},
		},
		{
			`
			<!doctype html>
			<body>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			</body>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{Name: "head"},
					ast.Element{Name: "body"},
				},
			},
		},
		{
			`
			<!doctype html>
			<html>
				<head>
					<meta charset="utf8">
					<title>Howdy!</title>
					<script src="main.js"></script>
				</head>
				<body>

				</body>
			</html>
			`,
			ast.Element{
				Name: "html",
				Children: []ast.Node{
					ast.Element{
						Name: "head",
						Children: []ast.Node{
							ast.Element{
								Name: "meta",
								Attributes: []ast.Attribute{
									ast.Attribute{Name: "charset", Value: "utf8"},
								},
							},
							ast.Element{
								Name: "title",
								Children: []ast.Node{
									ast.Text{Data: "Howdy!"},
								},
							},
							ast.Element{
								Name: "script",
								Attributes: []ast.Attribute{
									ast.Attribute{Name: "src", Value: "main.js"},
								},
							},
						},
					},
					ast.Element{Name: "body"},
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

		assert.Equal(t, test.root, ast, test.input)
	}
}
