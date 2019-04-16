package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/css/token"
)

func TestScan(t *testing.T) {
	var tests = []struct {
		input  string
		tokens []token.Token
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
				token.Delim{Offset: 0, Value: '#'},
				token.Ident{Offset: 1, Value: "foo"},
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
				token.Function{Offset: 0, Value: "url"},
				token.String{Offset: 4, Value: "foo"},
				token.CloseParen{Offset: 9},
			},
		},
		{
			"url('foo')",
			[]token.Token{
				token.Function{Offset: 0, Value: "url"},
				token.String{Offset: 4, Value: "foo"},
				token.CloseParen{Offset: 9},
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
			"foo(\"bar\")",
			[]token.Token{
				token.Function{Offset: 0, Value: "foo"},
				token.String{Offset: 4, Value: "bar"},
				token.CloseParen{Offset: 9},
			},
		},
		{
			"@foo",
			[]token.Token{
				token.AtKeyword{Offset: 0, Value: "foo"},
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
		{
			"123%",
			[]token.Token{
				token.Percentage{Offset: 0, Value: 1.23},
			},
		},
		{
			"123px",
			[]token.Token{
				token.Dimension{Offset: 0, Value: 123, Integer: true, Unit: "px"},
			},
		},
		{
			"1.23px",
			[]token.Token{
				token.Dimension{Offset: 0, Value: 1.23, Unit: "px"},
			},
		},
		{
			`\61\62\63`,
			[]token.Token{
				token.Ident{Offset: 0, Value: "abc"},
			},
		},
	}

	for _, test := range tests {
		tokens, err := Scan([]rune(test.input))

		assert.Nil(t, err)
		assert.Equal(t, test.tokens, tokens, test.input)
	}
}

func TestScanError(t *testing.T) {
	var tests = []struct {
		input  string
		tokens []token.Token
		err    error
	}{
		{
			"\"foo",
			[]token.Token{
				token.String{Offset: 0, Value: "foo"},
			},
			SyntaxError{Offset: 4, Message: "unexpected end of file"},
		},
		{
			"'foo",
			[]token.Token{
				token.String{Offset: 0, Value: "foo"},
			},
			SyntaxError{Offset: 4, Message: "unexpected end of file"},
		},
		{
			"\"foo\n",
			[]token.Token{},
			SyntaxError{Offset: 4, Message: "unexpected newline"},
		},
		{
			"'foo\n",
			[]token.Token{},
			SyntaxError{Offset: 4, Message: "unexpected newline"},
		},
		{
			"\\\n",
			[]token.Token{
				token.Delim{Offset: 0, Value: '\\'},
			},
			SyntaxError{Offset: 1, Message: "unexpected newline"},
		},
		{
			"url(foo",
			[]token.Token{
				token.Url{Offset: 0, Value: "foo"},
			},
			SyntaxError{Offset: 7, Message: "unexpected end of file"},
		},
		{
			"url(foo bar)",
			[]token.Token{},
			SyntaxError{Offset: 7, Message: "unexpected whitespace"},
		},
	}

	for _, test := range tests {
		tokens, err := Scan([]rune(test.input))

		assert.Equal(t, test.err, err, test.input)
		assert.Equal(t, test.tokens, tokens, test.input)
	}
}
