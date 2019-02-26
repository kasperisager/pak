package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/css/token"
)

func TestScan(t *testing.T) {
	var tests = []struct {
		input    string
		expected []token.Token
	}{
		{
			"--foo",
			[]token.Token{
				token.Ident{Offset: 0, Value: "--foo"},
			},
		},
		{
			"#foo",
			[]token.Token{
				token.Hash{Offset: 0, Value: "foo", Id: true},
			},
		},
		{
			"#123",
			[]token.Token{
				token.Hash{Offset: 0, Value: "123"},
			},
		},
		{
			".foo",
			[]token.Token{
				token.Delim{Offset: 0, Value: '.'},
				token.Ident{Offset: 1, Value: "foo"},
			},
		},
		{
			"\"foo\"",
			[]token.Token{
				token.String{Offset: 0, Value: "foo"},
			},
		},
		{
			"'foo'",
			[]token.Token{
				token.String{Offset: 0, Value: "foo"},
			},
		},
		{
			"url(foo)",
			[]token.Token{
				token.Url{Offset: 0, Value: "foo"},
			},
		},
		{
			"url(\"foo\")",
			[]token.Token{
				token.Url{Offset: 0, Value: "foo"},
			},
		},
		{
			"url('foo')",
			[]token.Token{
				token.Url{Offset: 0, Value: "foo"},
			},
		},
		{
			"foo(bar)",
			[]token.Token{
				token.Function{Offset: 0, Value: "foo"},
				token.Ident{Offset: 4, Value: "bar"},
				token.CloseParen{Offset: 7},
			},
		},
		{
			"123",
			[]token.Token{
				token.Number{Offset: 0, Value: 123, Integer: true},
			},
		},
		{
			"1.23",
			[]token.Token{
				token.Number{Offset: 0, Value: 1.23},
			},
		},
		{
			"1e2",
			[]token.Token{
				token.Number{Offset: 0, Value: 100},
			},
		},
		{
			"1e-2",
			[]token.Token{
				token.Number{Offset: 0, Value: 0.01},
			},
		},
		{
			"1.23e2",
			[]token.Token{
				token.Number{Offset: 0, Value: 123},
			},
		},
		{
			"1.23e-2",
			[]token.Token{
				token.Number{Offset: 0, Value: 0.0123},
			},
		},
	}

	for _, test := range tests {
		tokens, err := Scan([]rune(test.input))

		assert.Nil(t, err)
		assert.Equal(t, test.expected, tokens, test.input)
	}
}
