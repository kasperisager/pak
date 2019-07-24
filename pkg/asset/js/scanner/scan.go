package scanner

import (
	"unicode"

	"github.com/kasperisager/pak/pkg/asset/js/token"
	. "github.com/kasperisager/pak/pkg/runes"
)

type SyntaxError struct {
	Offset  int
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

type Options struct {
	regExp       bool
	templateTail bool
}

func RegExp(allowed bool) func(*Options) {
	return func(options *Options) {
		options.regExp = allowed
	}
}

func TemplateTail(allowed bool) func(*Options) {
	return func(options *Options) {
		options.templateTail = allowed
	}
}

func Scan(offset int, runes []rune, options ...func(*Options)) (int, []rune, token.Token) {
	_options := Options{
		regExp:       true,
		templateTail: false,
	}

	for _, option := range options {
		option(&_options)
	}

	return scanToken(offset, runes, _options)
}

func ScanInto(offset int, runes []rune, tokens []token.Token, options ...func(*Options)) (int, []rune, []token.Token, bool) {
	offset, runes, token := Scan(offset, runes, options...)

	if token == nil {
		return offset, runes, tokens, false
	}

	return offset, runes, append(tokens, token), true
}

func scanToken(offset int, runes Runes, options Options) (int, Runes, token.Token) {
	start := offset

	if offset, runes, value, ok := scanIdentifierName(offset, runes); ok {
		if isKeyword(value) || value == "enum" {
			return offset, runes, token.Keyword{Offset: start, Value: value}
		}

		switch value {
		case "true":
			return offset, runes, token.Boolean{Offset: start, Value: true}

		case "false":
			return offset, runes, token.Boolean{Offset: start, Value: false}

		case "null":
			return offset, runes, token.Null{Offset: start}
		}

		return offset, runes, token.Identifier{Offset: start, Value: value}
	}

	if offset, runes, ok := scanWhitespace(offset, runes); ok {
		return offset, runes, token.Whitespace{Offset: start}
	}

	if offset, runes, value, ok := scanPunctuator(offset, runes, options); ok {
		return offset, runes, token.Punctuator{Offset: start, Value: value}
	}

	if offset, runes, value, ok := scanStringLiteral(offset, runes); ok {
		return offset, runes, token.String{Offset: start, Value: value}
	}

	return offset, runes, nil
}

// https://www.ecma-international.org/ecma-262/#prod-UnicodeEscapeSequence
func scanUnicodeEscapeSequence(offset int, runes Runes) (int, Runes, rune, bool) {
	if runes.Peek(1) == 'u' {
		if runes.Peek(2) == '{' {
			offset, runes, rune, ok := scanCodePoint(offset+2, runes[2:])

			if ok && runes.Peek(1) == '}' {
				return offset, runes, rune, true
			}
		} else {
			for i := 0; i < 4; i++ {
			}
		}
	}

	return offset, runes, -1, false
}

// https://www.ecma-international.org/ecma-262/#prod-CodePoint
func scanCodePoint(offset int, runes Runes) (int, Runes, rune, bool) {
	var code int

	for i := 0; i < len(runes); i++ {
		next := runes.Peek(i + 1)

		if IsHexDigit(next) {
			code = 0x10*code + HexValue(next)
		} else if code > 0x10ffff {
			break
		} else {
			return offset + i, runes[i:], rune(code), true
		}
	}

	return offset, runes, -1, false
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierName
func scanIdentifierName(offset int, runes Runes) (int, Runes, string, bool) {
	var result []rune

	offset, runes, rune, ok := scanIdentifierStart(offset, runes)

	if ok {
		result = append(result, rune)
	} else {
		return offset, runes, "", false
	}

	for {
		offset, runes, rune, ok = scanIdentifierStart(offset, runes)

		if ok {
			result = append(result, rune)
		} else {
			break
		}
	}

	for {
		offset, runes, rune, ok = scanIdentifierPart(offset, runes)

		if ok {
			result = append(result, rune)
		} else {
			break
		}
	}

	return offset, runes, string(result), true
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierStart
func scanIdentifierStart(offset int, runes Runes) (int, Runes, rune, bool) {
	next := runes.Peek(1)

	if isUnicodeIdentifierStart(next) || next == '$' || next == '_' {
		return offset + 1, runes[1:], next, true
	}

	if next == '\\' {
		offset, runes, code, ok := scanUnicodeEscapeSequence(offset+1, runes[1:])

		if ok {
			return offset, runes, code, true
		}
	}

	return offset, runes, -1, false
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierPart
func scanIdentifierPart(offset int, runes Runes) (int, Runes, rune, bool) {
	next := runes.Peek(1)

	if isUnicodeIdentifierContinue(next) || next == '$' || next == 0x200c || next == 0x200d {
		return offset + 1, runes[1:], next, true
	}

	if next == '\\' {
		offset, runes, code, ok := scanUnicodeEscapeSequence(offset+1, runes[1:])

		if ok {
			return offset, runes, code, true
		}
	}

	return offset, runes, -1, false
}

// https://www.ecma-international.org/ecma-262/#prod-HexDigit
func scanHexDigit(offset int, runes Runes) (int, Runes, int, bool) {
	next := runes.Peek(1)

	if IsHexDigit(next) {
		return offset + 1, runes[1:], HexValue(next), true
	}

	return offset, runes, -1, false
}

// https://www.ecma-international.org/ecma-262/#prod-WhiteSpace
func scanWhitespace(offset int, runes Runes) (int, Runes, bool) {
	next := runes.Peek(1)

	switch next {
	case ' ', '\t', 0xb, 0xc, 0xa0, 0xfeff:
		return offset + 1, runes[1:], true
	}

	if unicode.In(next, unicode.Zs) {
		return offset + 1, runes[1:], true
	}

	return offset, runes, false
}

// https://www.ecma-international.org/ecma-262/#prod-Punctuator
func scanPunctuator(offset int, runes Runes, options Options) (int, Runes, string, bool) {
	next := runes.Peek(1)

	switch next {
	case '{', '(', ')', '[', ']', ';', ',', '~', '?', ':':
		return offset + 1, runes[1:], string(next), true

	case '.':
		if runes.Peek(2) == '.' && runes.Peek(3) == '.' {
			return offset + 3, runes[3:], "...", true
		}

		return offset + 1, runes[1:], ".", true

	case '}':
		if !options.templateTail {
			return offset + 1, runes[1:], "}", true
		}

	case '/':
		if !options.regExp {
			if runes.Peek(2) == '=' {
				return offset + 2, runes[2:], "/=", true
			}

			return offset + 1, runes[1:], "/", true
		}
	}

	return offset, runes, "", false
}

// https://www.ecma-international.org/ecma-262/#prod-StringLiteral
func scanStringLiteral(offset int, runes Runes) (int, Runes, string, bool) {
	var mark rune

	switch next := runes.Peek(1); next {
	case '"', '\'':
		mark = next

	default:
		return offset, runes, "", false
	}

	for i := 1; i < len(runes); i++ {
		switch next := runes.Peek(i + 1); next {
		case mark:
			return offset + i + 1, runes[i+1:], string(runes[1:i]), true

		case '\n', '\r':
			break

		case '\\':
		}
	}

	return offset, runes, "", false
}

// https://www.ecma-international.org/ecma-262/#prod-LineTerminator
func isLineTerminator(rune rune) bool {
	switch rune {
	case '\n', '\r', 0x2028, 0x2029:
		return true

	default:
		return false
	}
}

// https://www.ecma-international.org/ecma-262/#prod-Keyword
func isKeyword(identifierName string) bool {
	switch identifierName {
	case
		"await",
		"break",
		"case",
		"catch",
		"class",
		"const",
		"continue",
		"debugger",
		"default",
		"delete",
		"do",
		"else",
		"export",
		"extends",
		"finally",
		"for",
		"function",
		"if",
		"import",
		"in",
		"instanceof",
		"new",
		"return",
		"super",
		"this",
		"throw",
		"try",
		"typeof",
		"var",
		"void",
		"while",
		"with",
		"yield":
		return true

	default:
		return false
	}
}

// https://unicode.org/reports/tr31/#Table_Lexical_Classes_for_Identifiers
func isUnicodeIdentifierStart(rune rune) bool {
	return unicode.In(rune,
		unicode.Lu,
		unicode.Ll,
		unicode.Lt,
		unicode.Lm,
		unicode.Lo,
		unicode.Nl,
		unicode.Other_ID_Start,
	) && !unicode.In(rune,
		unicode.Pattern_Syntax,
		unicode.Pattern_White_Space,
	)
}

// https://unicode.org/reports/tr31/#Table_Lexical_Classes_for_Identifiers
func isUnicodeIdentifierContinue(rune rune) bool {
	return isUnicodeIdentifierStart(rune) || unicode.In(rune,
		unicode.Mn,
		unicode.Mc,
		unicode.Nd,
		unicode.Pc,
		unicode.Other_ID_Continue,
	) && !unicode.In(rune,
		unicode.Pattern_Syntax,
		unicode.Pattern_White_Space,
	)
}
