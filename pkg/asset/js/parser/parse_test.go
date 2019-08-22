package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/js/ast"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input string
		program *ast.Program
	}{
		{
			`foo = "bar"`,
			&ast.Program{
				Body: []ast.ProgramBody{
					&ast.ExpressionStatement{
						Expression: &ast.AssignmentExpression{
							Operator: "=",
							Left: &ast.Identifier{Name: "foo"},
							Right: &ast.StringLiteral{Value: "bar"},
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
							Left: &ast.Identifier{Name: "foo"},
							Right: &ast.StringLiteral{Value: "bar"},
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

