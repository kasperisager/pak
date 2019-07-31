package scanner

import (
	"unicode"

	"github.com/kasperisager/pak/pkg/asset/js/token"
	"github.com/kasperisager/pak/pkg/runes"
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

func Scan(offset int, runes []rune, options ...func(*Options)) (int, []rune, token.Token, error) {
	_options := Options{
		regExp:       true,
		templateTail: false,
	}

	for _, option := range options {
		option(&_options)
	}

	scanner, token, err := scanToken(scanner{offset, runes}, _options)

	return scanner.offset, scanner.runes, token, err
}

type (
	scanner struct {
		offset int
		runes  []rune
	}
)

const eof = -1

func (s scanner) peek(n int) (scanner, rune) {
	if len(s.runes) >= n {
		return s, s.runes[n-1]
	}

	return s, eof
}

func (s scanner) advance(n int) scanner {
	if len(s.runes) >= n {
		s.runes = s.runes[n:]
	}

	return s
}

func scanToken(scanner scanner, options Options) (scanner, token.Token, error) {
	for {
		var ok bool

		scanner, ok = scanWhitespace(scanner)

		if !ok {
			break
		}
	}

	start := scanner.offset

	scanner, next := scanner.peek(1)

	switch next {
	case '{', '(', ')', '[', ']', ';', ',', '~', '?', ':', '.', '=':
		scanner, value, err := scanPunctuator(scanner, options)

		if err != nil {
			return scanner, nil, err
		}

		return scanner, token.Punctuator{Offset: start, Value: value}, nil

	case '"', '\'':
		scanner, value, err := scanStringLiteral(scanner)

		if err != nil {
			return scanner, nil, err
		}

		return scanner, token.String{Offset: start, Value: value}, nil

	default:
		if scanner, value, err := scanIdentifierName(scanner); err == nil {
			if isKeyword(value) || value == "enum" {
				return scanner, token.Keyword{Offset: start, Value: value}, nil
			}

			switch value {
			case "true":
				return scanner, token.Boolean{Offset: start, Value: true}, nil

			case "false":
				return scanner, token.Boolean{Offset: start, Value: false}, nil

			case "null":
				return scanner, token.Null{Offset: start}, nil
			}

			return scanner, token.Identifier{Offset: start, Value: value}, nil
		}

		return scanner, nil, SyntaxError{
			Offset:  start,
			Message: "unexpected character",
		}
	}
}

// https://www.ecma-international.org/ecma-262/#prod-UnicodeEscapeSequence
func scanUnicodeEscapeSequence(scanner scanner) (scanner, rune, error) {
	var next rune

	if scanner, next = scanner.peek(1); next == 'u' {
		if scanner, next = scanner.peek(2); next == '{' {
			scanner, codePoint, err := scanCodePoint(scanner.advance(2))

			if err != nil {
				return scanner, -1, err
			}

			scanner, next := scanner.peek(1)

			if next != '}' {
				return scanner, -1, SyntaxError{
					Offset:  scanner.offset,
					Message: `unexpected character, expected "}"`,
				}
			}

			return scanner, codePoint, nil
		} else {
			for i := 0; i < 4; i++ {
			}
		}
	}

	return scanner, -1, SyntaxError{
		Offset:  scanner.offset,
		Message: "unexpected character, expected escape sequence",
	}
}

// https://www.ecma-international.org/ecma-262/#prod-CodePoint
func scanCodePoint(scanner scanner) (scanner, rune, error) {
	var codePoint int

	for {
		var next rune

		scanner, next = scanner.peek(1)

		if runes.IsHexDigit(next) {
			codePoint = 0x10*codePoint + runes.HexValue(next)
		} else if codePoint <= 0x10ffff {
			break
		} else {
			return scanner, -1, SyntaxError{
				Offset:  scanner.offset,
				Message: "unexpected code point larger than 0x10ffff",
			}
		}
	}

	return scanner, rune(codePoint), nil
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierName
func scanIdentifierName(scanner scanner) (scanner, string, error) {
	var result []rune

	scanner, rune, err := scanIdentifierStart(scanner)

	if err != nil {
		return scanner, "", err
	}

	result = append(result, rune)

	for {
		scanner, rune, err = scanIdentifierStart(scanner)

		if err != nil {
			break
		}

		result = append(result, rune)
	}

	for {
		scanner, rune, err = scanIdentifierPart(scanner)

		if err != nil {
			break
		}

		result = append(result, rune)
	}

	return scanner, string(result), nil
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierStart
func scanIdentifierStart(scanner scanner) (scanner, rune, error) {
	scanner, next := scanner.peek(1)

	if isUnicodeIdentifierStart(next) || next == '$' || next == '_' {
		return scanner.advance(1), next, nil
	}

	if next == '\\' {
		scanner, escape, err := scanUnicodeEscapeSequence(scanner.advance(1))

		if err != nil {
			return scanner, -1, err
		}

		return scanner, escape, nil
	}

	return scanner, -1, SyntaxError{
		Offset:  scanner.offset,
		Message: "unexpected character, expected start of identifier",
	}
}

// https://www.ecma-international.org/ecma-262/#prod-IdentifierPart
func scanIdentifierPart(scanner scanner) (scanner, rune, error) {
	scanner, next := scanner.peek(1)

	if isUnicodeIdentifierContinue(next) || next == '$' || next == 0x200c || next == 0x200d {
		return scanner.advance(1), next, nil
	}

	if next == '\\' {
		scanner, code, err := scanUnicodeEscapeSequence(scanner.advance(1))

		if err != nil {
			return scanner, -1, err
		}

		return scanner, code, nil
	}

	return scanner, -1, SyntaxError{
		Offset:  scanner.offset,
		Message: "unexpected character, expected part of identifier",
	}
}

// https://www.ecma-international.org/ecma-262/#prod-HexDigit
func scanHexDigit(scanner scanner) (scanner, int, bool) {
	scanner, next := scanner.peek(1)

	if runes.IsHexDigit(next) {
		return scanner.advance(1), runes.HexValue(next), true
	}

	return scanner, -1, false
}

// https://www.ecma-international.org/ecma-262/#prod-WhiteSpace
func scanWhitespace(scanner scanner) (scanner, bool) {
	scanner, next := scanner.peek(1)

	switch next {
	case ' ', '\t', 0xb, 0xc, 0xa0, 0xfeff:
		return scanner.advance(1), true
	}

	if unicode.In(next, unicode.Zs) {
		return scanner.advance(1), true
	}

	return scanner, false
}

// https://www.ecma-international.org/ecma-262/#prod-Punctuator
func scanPunctuator(scanner scanner, options Options) (scanner, string, error) {
	scanner, next := scanner.peek(1)

	switch next {
	case '{', '(', ')', '[', ']', ';', ',', '~', '?', ':':
		return scanner.advance(1), string(next), nil

	case '.':
		scanner, next = scanner.peek(2)

		if next == '.' {
			scanner, next = scanner.peek(3)

			if next == '.' {
				return scanner.advance(3), "...", nil
			}
		}

		return scanner.advance(1), ".", nil

	case '}':
		if !options.templateTail {
			return scanner.advance(1), "}", nil
		}

	case '/':
		if !options.regExp {
			scanner, next = scanner.peek(2)

			if next == '=' {
				return scanner.advance(2), "/=", nil
			}

			return scanner.advance(1), "/", nil
		}

	case '=':
		return scanner.advance(1), "=", nil
	}

	return scanner, "", SyntaxError{
		Offset:  scanner.offset,
		Message: "unexpected character, expected punctuator",
	}
}

// https://www.ecma-international.org/ecma-262/#prod-StringLiteral
func scanStringLiteral(scanner scanner) (scanner, string, error) {
	var mark rune

	scanner, next := scanner.peek(1)

	switch next {
	case '"', '\'':
		mark = next
		scanner = scanner.advance(1)

	default:
		return scanner, "", SyntaxError{
			Offset:  scanner.offset,
			Message: `unexpected character, expected start of string`,
		}
	}

	var value []rune

	for {
		scanner, next = scanner.peek(1)

		switch next {
		default:
			value = append(value, next)
			scanner = scanner.advance(1)

		case mark:
			return scanner.advance(1), string(value), nil

		case '\n', '\r':
			return scanner, "", SyntaxError{
				Offset:  scanner.offset,
				Message: "unexpected newline in string literal",
			}

		case eof:
			return scanner, "", SyntaxError{
				Offset:  scanner.offset,
				Message: "unexpected unterminated string literal",
			}
		}
	}
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
