package scanner

import (
	"strings"
	"math"

	"github.com/kasperisager/pak/pkg/asset/css/token"
)

func Scan(runes []rune) (tokens []token.Token, err error) {
	offset := 0

	for len(runes) > 0 {
		offset, runes, tokens, err = scanToken(offset, runes, tokens)

		if err != nil {
			return tokens, err
		}
	}

	return tokens, nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-token
func scanToken(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	switch runes[0] {
	case '+', '.':
		if startsNumber(runes) {
			return scanNumeric(offset, runes, tokens)
		}

	case '-':
		if startsNumber(runes) {
			return scanNumeric(offset, runes, tokens)
		}

		if startsIdentifier(runes) {
			return scanIdent(offset, runes, tokens)
		}

	case '\\':
		if startsEscape(runes) {
			return scanIdent(offset, runes, tokens)
		}

		return offset, runes, tokens, nil

	case '"', '\'':
		return scanString(offset, runes, tokens)

	case '#':
		if isNameStart(runes[1]) || startsEscape(runes[1:]) {
			start := offset

			offset, runes, name, err := scanName(offset+1, runes[1:])

			if err != nil {
				return offset, runes, tokens, err
			}

			t := token.Hash{
				Offset: start,
				Value:  name,
			}

			return offset, runes, append(tokens, t), nil
		}

	case '@':
		if startsIdentifier(runes[1:]) {
			start := offset

			offset, runes, name, err := scanName(offset+1, runes[1:])

			if err != nil {
				return offset, runes, tokens, err
			}

			t := token.AtKeyword{
				Offset: start,
				Value:  name,
			}

			return offset, runes, append(tokens, t), nil
		}

	case '/':
		if len(runes) > 1 && runes[1] == '*' {
			return scanComment(offset+2, runes[2:], tokens)
		}

	case ',':
		return offset + 1, runes[1:], append(tokens, token.Comma{Offset: offset}), nil

	case ':':
		return offset + 1, runes[1:], append(tokens, token.Colon{Offset: offset}), nil

	case ';':
		return offset + 1, runes[1:], append(tokens, token.Semicolon{Offset: offset}), nil

	case '(':
		return offset + 1, runes[1:], append(tokens, token.OpenParen{Offset: offset}), nil

	case ')':
		return offset + 1, runes[1:], append(tokens, token.CloseParen{Offset: offset}), nil

	case '[':
		return offset + 1, runes[1:], append(tokens, token.OpenSquare{Offset: offset}), nil

	case ']':
		return offset + 1, runes[1:], append(tokens, token.CloseSquare{Offset: offset}), nil

	case '{':
		return offset + 1, runes[1:], append(tokens, token.OpenCurly{Offset: offset}), nil

	case '}':
		return offset + 1, runes[1:], append(tokens, token.CloseCurly{Offset: offset}), nil

	default:
		switch {
		case isDigit(runes[0]):
			return scanNumeric(offset, runes, tokens)

		case isNameStart(runes[0]):
			return scanIdent(offset, runes, tokens)

		case isWhitespace(runes[0]):
			return scanWhitespace(offset, runes, tokens)
		}
	}

	return offset + 1, runes[1:], append(tokens, token.Delim{Offset: offset, Value: runes[0]}), nil
}

func scanWhitespace(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	start := offset

	for len(runes) > 0 && isWhitespace(runes[0]) {
		runes = runes[1:]
		offset++
	}

	t := token.Whitespace{
		Offset: start,
	}

	return offset, runes, append(tokens, t), nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-comments
func scanComment(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	for len(runes) > 0 {
		if runes[0] == '*' {
			runes = runes[1:]
			offset++

			if len(runes) > 0 && runes[0] == '/' {
				runes = runes[1:]
				offset++
				break
			}
		} else {
			runes = runes[1:]
			offset++
		}
	}

	return offset, runes, tokens, nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-numeric-token
func scanNumeric(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	offset, runes, value, integer := scanNumber(offset, runes)

	t := token.Number{
		Value:   value,
		Integer: integer,
	}

	return offset, runes, append(tokens, t), nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-number
func scanNumber(offset int, runes []rune) (int, []rune, float64, bool) {
	value := 0.0
	sign := 1.0

	switch runes[0] {
	case '+':
		runes = runes[1:]
		offset++

	case '-':
		sign = -1
		runes = runes[1:]
		offset++
	}

	for len(runes) > 0 {
		if isDigit(runes[0]) {
			value = 10*value + float64(runes[0] - '0')
			runes = runes[1:]
			offset++
		} else {
			break
		}
	}

	if startsFraction(runes) {
		offset, runes, value = scanFraction(offset + 1, runes[1:], value)

		return offset, runes, sign * value, false
	}

	if startsExponent(runes) {
		offset, runes, value = scanExponent(offset + 1, runes[1:], value)

		return offset, runes, sign * value, false
	}

	return offset, runes, sign * value, true
}

func scanFraction(offset int, runes []rune, base float64) (int, []rune, float64) {
	value := 0.0
	digits := 0

	for len(runes) > 0 {
		if isDigit(runes[0]) {
			value = 10*value + float64(runes[0] - '0')
			digits++
			runes = runes[1:]
			offset++
		} else {
			break
		}
	}

	value = base + value / math.Pow10(digits)

	if startsExponent(runes) {
		return scanExponent(offset + 1, runes[1:], value)
	}

	return offset, runes, value
}

func scanExponent(offset int, runes []rune, base float64) (int, []rune, float64) {
	value := 0
	sign := -1

	switch runes[0] {
	case '+':
		runes = runes[1:]
		offset++
	case '-':
		runes = runes[1:]
		offset++
		sign = 1
	}

	for len(runes) > 0 {
		if isDigit(runes[0]) {
			value = 10*value + int(runes[0] - '0')
			runes = runes[1:]
			offset++
		} else {
			break
		}
	}

	return offset, runes, base / math.Pow10(sign * value)
}

// See: https://drafts.csswg.org/css-syntax/#consume-a-string-token
func scanString(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	start := offset
	end := runes[0]
	runes = runes[1:]
	offset++

	var result strings.Builder

	for len(runes) > 0 {
		if isNewline(runes[0]) {
			return offset, runes, tokens, nil
		}

		if runes[0] == end {
			t := token.String{
				Offset: start,
				Value:  result.String(),
			}

			return offset + 1, runes[1:], append(tokens, t), nil
		}

		result.WriteRune(runes[0])

		runes = runes[1:]
		offset++
	}

	return offset, runes, tokens, nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-a-name
func scanName(offset int, runes []rune) (int, []rune, string, error) {
	var result strings.Builder

	for len(runes) > 0 {
		switch {
		case isName(runes[0]):
			result.WriteRune(runes[0])
			runes = runes[1:]
			offset++

		case startsEscape(runes):
			var (
				code rune
				err  error
			)

			offset, runes, code, err = scanEscape(offset, runes)

			if err != nil {
				return offset, runes, result.String(), err
			}

			result.WriteRune(code)

		default:
			return offset, runes, result.String(), nil
		}
	}

	return offset, runes, result.String(), nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-escaped-code-point
func scanEscape(offset int, runes []rune) (int, []rune, rune, error) {
	if isHexDigit(runes[0]) {
		code := hexValue(runes[0])

		for i := 0; len(runes) > 0 && i < 5; i++ {
			if isHexDigit(runes[0]) {
				code = 0x10*code + hexValue(runes[0])
				runes = runes[1:]
				offset++
			} else {
				break
			}
		}

		if len(runes) > 0 && isWhitespace(runes[0]) {
			runes = runes[1:]
			offset++
		}

		return offset, runes, rune(code), nil
	}

	return offset + 1, runes[1:], runes[0], nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-ident-like-token
func scanIdent(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	start := offset

	var (
		name string
		err  error
	)

	offset, runes, name, err = scanName(offset, runes)

	if err != nil {
		return offset, runes, tokens, err
	}

	var t token.Token

	if len(runes) > 0 && runes[0] == '(' {
		runes = runes[1:]
		offset++

		if strings.EqualFold(name, "url") {
			return scanUrl(offset, runes, tokens)
		}

		t = token.Function{
			Offset: start,
			Value:  name,
		}
	} else {
		t = token.Ident{
			Offset: start,
			Value:  name,
		}
	}

	return offset, runes, append(tokens, t), nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-url-token
func scanUrl(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	for len(runes) > 0 && isWhitespace(runes[0]) {
		runes = runes[1:]
		offset++
	}

	var result strings.Builder

	for len(runes) > 0 {
		switch runes[0] {
		case ')':
			t := token.Url{
				Value: result.String(),
			}

			return offset + 1, runes[1:], append(tokens, t), nil

		case '\\':
			if startsEscape(runes) {
				var (
					code rune
					err  error
				)

				offset, runes, code, err = scanEscape(offset, runes)

				if err != nil {
					return offset, runes, tokens, err
				}

				result.WriteRune(code)
			} else {
				result.WriteRune(runes[0])
				runes = runes[1:]
				offset++
			}

		case '"', '\'':
			end := runes[0]
			runes = runes[1:]
			offset++

			for len(runes) > 0 {
				if runes[0] == end {
					runes = runes[1:]
					offset++
					break
				}

				result.WriteRune(runes[0])
				runes = runes[1:]
				offset++
			}

		default:
			result.WriteRune(runes[0])
			runes = runes[1:]
			offset++
		}
	}

	return offset, runes, tokens, nil
}
