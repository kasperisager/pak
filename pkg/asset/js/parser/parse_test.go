package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/js/ast"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input   string
		program *ast.Program
	}{
		{
			`foo = "bar"`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.AssignmentExpression{
							Operator: "=",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.StringLiteral{Value: "bar"},
						},
					},
				},
			},
		},
		{
			`foo, bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.SequenceExpression{
							Expression: []ast.Expression{
								&ast.Identifier{Name: "foo"},
								&ast.Identifier{Name: "bar"},
							},
						},
					},
				},
			},
		},
		{
			`foo *= "bar"`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.AssignmentExpression{
							Operator: "*=",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.StringLiteral{Value: "bar"},
						},
					},
				},
			},
		},
		{
			`foo ? bar : baz`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.ConditionalExpression{
							Test:       &ast.Identifier{Name: "foo"},
							Alternate:  &ast.Identifier{Name: "bar"},
							Consequent: &ast.Identifier{Name: "baz"},
						},
					},
				},
			},
		},
		{
			`foo || bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.LogicalExpression{
							Operator: "||",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo && bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.LogicalExpression{
							Operator: "&&",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo || bar || baz`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.LogicalExpression{
							Operator: "||",
							Left:     &ast.Identifier{Name: "foo"},
							Right: &ast.LogicalExpression{
								Operator: "||",
								Left:     &ast.Identifier{Name: "bar"},
								Right:    &ast.Identifier{Name: "baz"},
							},
						},
					},
				},
			},
		},
		{
			`foo && bar && baz`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.LogicalExpression{
							Operator: "&&",
							Left:     &ast.Identifier{Name: "foo"},
							Right: &ast.LogicalExpression{
								Operator: "&&",
								Left:     &ast.Identifier{Name: "bar"},
								Right:    &ast.Identifier{Name: "baz"},
							},
						},
					},
				},
			},
		},
		{
			`foo || bar && baz`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.LogicalExpression{
							Operator: "||",
							Left:     &ast.Identifier{Name: "foo"},
							Right: &ast.LogicalExpression{
								Operator: "&&",
								Left:     &ast.Identifier{Name: "bar"},
								Right:    &ast.Identifier{Name: "baz"},
							},
						},
					},
				},
			},
		},
		{
			`foo && bar || baz`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.LogicalExpression{
							Operator: "||",
							Left: &ast.LogicalExpression{
								Operator: "&&",
								Left:     &ast.Identifier{Name: "foo"},
								Right:    &ast.Identifier{Name: "bar"},
							},
							Right: &ast.Identifier{Name: "baz"},
						},
					},
				},
			},
		},
		{
			`foo | bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "|",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo ^ bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "^",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo & bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "&",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo == bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "==",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo != bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "!=",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo === bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "===",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo !== bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "!==",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo < bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "<",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo > bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: ">",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo <= bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "<=",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo >= bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: ">=",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo << bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "<<",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo >> bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: ">>",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo >>> bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: ">>>",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo + bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "+",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`foo - bar`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Operator: "-",
							Left:     &ast.Identifier{Name: "foo"},
							Right:    &ast.Identifier{Name: "bar"},
						},
					},
				},
			},
		},
		{
			`+foo`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "+",
							Prefix:   true,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`-foo`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "-",
							Prefix:   true,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`~foo`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "~",
							Prefix:   true,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`!foo`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Prefix:   true,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`++foo`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UpdateExpression{
							Operator: "++",
							Prefix:   true,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`--foo`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UpdateExpression{
							Operator: "--",
							Prefix:   true,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`foo++`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UpdateExpression{
							Operator: "++",
							Prefix:   false,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
		{
			`foo--`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.UpdateExpression{
							Operator: "--",
							Prefix:   false,
							Argument: &ast.Identifier{Name: "foo"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		program, err := Parse([]rune(test.input))
		assert.Nil(t, err, test.input)

		assert.Equal(t, test.program, program, test.input)
	}
}
