package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/js/token"
)

func TestScan(t *testing.T) {
	var tests = []struct {
		input  string
		tokens []token.Token
	}{
		{
			`foo`,
			[]token.Token{
				token.Identifier{Offset: 0, Value: "foo"},
			},
		},
		{
			`"foo"`,
			[]token.Token{
				token.String{Offset: 0, Value: "foo"},
			},
		},
		{
			`'foo'`,
			[]token.Token{
				token.String{Offset: 0, Value: "foo"},
			},
		},
		{
			`123`,
			[]token.Token{
				token.Number{Offset: 0, Value: 123},
			},
		},
		{
			`1.23`,
			[]token.Token{
				token.Number{Offset: 0, Value: 1.23},
			},
		},
		{
			`.123`,
			[]token.Token{
				token.Number{Offset: 0, Value: .123},
			},
		},
		{
			`1.23e4`,
			[]token.Token{
				token.Number{Offset: 0, Value: 1.23e4},
			},
		},
		{
			`1.23e+4`,
			[]token.Token{
				token.Number{Offset: 0, Value: 1.23e+4},
			},
		},
		{
			`1.23e-4`,
			[]token.Token{
				token.Number{Offset: 0, Value: 1.23e-4},
			},
		},
	}

	for _, test := range tests {
		runes := []rune(test.input)

		var (
			offset int
			next   token.Token
			err    error
			tokens []token.Token
		)

		for {
			offset, runes, next, err = Scan(offset, runes)

			if err != nil {
				break
			}

			tokens = append(tokens, next)
		}

		assert.Equal(t, test.tokens, tokens, test.input)
	}
}
