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
		{
			`\u{66}\u{6f}\u{6f}`,
			[]token.Token{
				token.Identifier{Offset: 0, Value: "foo"},
			},
		},
		{
			`\u0066\u006f\u006f`,
			[]token.Token{
				token.Identifier{Offset: 0, Value: "foo"},
			},
		},
		{
			`||`,
			[]token.Token{
				token.Punctuator{Offset: 0, Value: "||"},
			},
		},
		{
			`&&`,
			[]token.Token{
				token.Punctuator{Offset: 0, Value: "&&"},
			},
		},
		{
			`!`,
			[]token.Token{
				token.Punctuator{Offset: 0, Value: "!"},
			},
		},
		{
			`!foo`,
			[]token.Token{
				token.Punctuator{Offset: 0, Value: "!"},
				token.Identifier{Offset: 1, Value: "foo"},
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
