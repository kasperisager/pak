package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kasperisager/pak/pkg/asset/html/token"
)

func TestScan(t *testing.T) {
	var tests = []struct {
		input  string
		tokens []token.Token
	}{
		{
			`<!doctype html>`,
			[]token.Token{
				token.DocumentType{Offset: 0},
			},
		},
		{
			`<foo>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo"},
			},
		},
		{
			`</foo>`,
			[]token.Token{
				token.EndTag{Offset: 0, Name: "foo"},
			},
		},
		{
			`<foo/>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo", Closed: true},
			},
		},
		{
			`<foo></foo>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo"},
				token.EndTag{Offset: 5, Name: "foo"},
			},
		},
		{
			`<foo bar>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo", Attributes: []token.Attribute{
					token.Attribute{Offset: 5, Name: "bar"},
				}},
			},
		},
		{
			`<foo bar=baz>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo", Attributes: []token.Attribute{
					token.Attribute{Offset: 5, Name: "bar", Value: "baz"},
				}},
			},
		},
		{
			`<foo bar="baz">`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo", Attributes: []token.Attribute{
					token.Attribute{Offset: 5, Name: "bar", Value: "baz"},
				}},
			},
		},
		{
			`<foo bar baz>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo", Attributes: []token.Attribute{
					token.Attribute{Offset: 5, Name: "bar"},
					token.Attribute{Offset: 9, Name: "baz"},
				}},
			},
		},
		{
			`<foo bar=baz qux>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "foo", Attributes: []token.Attribute{
					token.Attribute{Offset: 5, Name: "bar", Value: "baz"},
					token.Attribute{Offset: 13, Name: "qux"},
				}},
			},
		},
		{
			`<script><foo></script>`,
			[]token.Token{
				token.StartTag{Offset: 0, Name: "script"},
				token.Character{Offset: 8, Data: '<'},
				token.Character{Offset: 9, Data: 'f'},
				token.Character{Offset: 10, Data: 'o'},
				token.Character{Offset: 11, Data: 'o'},
				token.Character{Offset: 12, Data: '>'},
				token.EndTag{Offset: 13, Name: "script"},
			},
		},
		{
			`<!--foo-->`,
			[]token.Token{},
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
	}{}

	for _, test := range tests {
		tokens, err := Scan([]rune(test.input))

		assert.Equal(t, test.err, err, test.input)
		assert.Equal(t, test.tokens, tokens, test.input)
	}
}
