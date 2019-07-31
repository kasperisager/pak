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
			`foo = "foo"`,
			&ast.Program{},
		},
	}

	for _, test := range tests {
		program, err := Parse([]rune(test.input))
		assert.Nil(t, err, test.input)

		assert.Equal(t, test.program, program, test.input)
	}
}

