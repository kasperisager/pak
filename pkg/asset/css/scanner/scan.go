package scanner

import (
	"math"
	"strings"

	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type SyntaxError struct {
	Offset  int
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

func Scan(runes []rune) (tokens []token.Token, err error) {
	tokens = make([]token.Token, 0, len(runes)/4)

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
	switch peek(runes, 1) {
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

		t := token.Delim{Offset: offset, Value: '\\'}

		return offset + 1, runes[1:], append(tokens, t), SyntaxError{
			Offset:  offset + 1,
			Message: "unexpected newline",
		}

	case '"', '\'':
		return scanString(offset, runes, tokens)

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
		if peek(runes, 2) == '*' {
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
		case isDigit(peek(runes, 1)):
			return scanNumeric(offset, runes, tokens)

		case isNameStart(peek(runes, 1)):
			return scanIdent(offset, runes, tokens)

		case isWhitespace(peek(runes, 1)):
			return scanWhitespace(offset, runes, tokens)
		}
	}

	return offset + 1, runes[1:], append(tokens, token.Delim{Offset: offset, Value: peek(runes, 1)}), nil
}

func scanWhitespace(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	start := offset

	for isWhitespace(peek(runes, 1)) {
		offset, runes = offset+1, runes[1:]
	}

	t := token.Whitespace{
		Offset: start,
	}

	return offset, runes, append(tokens, t), nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-comments
func scanComment(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	for len(runes) > 0 {
		if peek(runes, 1) == '*' && peek(runes, 2) == '/' {
			return offset + 2, runes[2:], tokens, nil
		}

		offset, runes = offset+1, runes[1:]
	}

	return offset, runes, tokens, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

// See: https://drafts.csswg.org/css-syntax/#consume-numeric-token
func scanNumeric(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	start := offset

	offset, runes, value, integer := scanNumber(offset, runes)

	var t token.Token

	switch {
	case startsIdentifier(runes):
		var (
			name string
			err  error
		)

		offset, runes, name, err = scanName(offset, runes)

		if err != nil {
			return offset, runes, tokens, err
		}

		t = token.Dimension{
			Offset:  start,
			Value:   value,
			Integer: integer,
			Unit:    name,
		}

	case peek(runes, 1) == '%':
		offset, runes = offset+1, runes[1:]

		t = token.Percentage{
			Offset: start,
			Value:  value / 100,
		}

	default:
		t = token.Number{
			Offset:  start,
			Value:   value,
			Integer: integer,
		}
	}

	return offset, runes, append(tokens, t), nil
}

// See: https://drafts.csswg.org/css-syntax/#consume-number
func scanNumber(offset int, runes []rune) (int, []rune, float64, bool) {
	value := 0.0
	sign := 1.0

	switch peek(runes, 1) {
	case '+':
		offset, runes = offset+1, runes[1:]

	case '-':
		sign = -1
		offset, runes = offset+1, runes[1:]
	}

	for len(runes) > 0 {
		if isDigit(peek(runes, 1)) {
			value = 10*value + float64(peek(runes, 1)-'0')
			offset, runes = offset+1, runes[1:]
		} else {
			break
		}
	}

	if startsFraction(runes) {
		offset, runes, value = scanFraction(offset+1, runes[1:], value)

		return offset, runes, sign * value, false
	}

	if startsExponent(runes) {
		offset, runes, value = scanExponent(offset+1, runes[1:], value)

		return offset, runes, sign * value, false
	}

	return offset, runes, sign * value, true
}

func scanFraction(offset int, runes []rune, base float64) (int, []rune, float64) {
	value := 0.0
	digits := 0

	for len(runes) > 0 {
		if isDigit(peek(runes, 1)) {
			value = 10*value + float64(peek(runes, 1)-'0')
			digits++
			offset, runes = offset+1, runes[1:]
		} else {
			break
		}
	}

	value = base + value/math.Pow10(digits)

	if startsExponent(runes) {
		return scanExponent(offset+1, runes[1:], value)
	}

	return offset, runes, value
}

func scanExponent(offset int, runes []rune, base float64) (int, []rune, float64) {
	value := 0
	sign := -1

	switch peek(runes, 1) {
	case '+':
		offset, runes = offset+1, runes[1:]

	case '-':
		sign = 1
		offset, runes = offset+1, runes[1:]
	}

	for len(runes) > 0 {
		if isDigit(peek(runes, 1)) {
			value = 10*value + int(peek(runes, 1)-'0')
			offset, runes = offset+1, runes[1:]
		} else {
			break
		}
	}

	return offset, runes, base / math.Pow10(sign*value)
}

// See: https://drafts.csswg.org/css-syntax/#consume-a-string-token
func scanString(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	start := offset
	end := peek(runes, 1)
	offset, runes = offset+1, runes[1:]

	var result strings.Builder

	for len(runes) > 0 {
		next := peek(runes, 1)

		if isNewline(next) {
			return offset, runes, tokens, SyntaxError{
				Offset:  offset,
				Message: "unexpected newline",
			}
		}

		offset, runes = offset+1, runes[1:]

		switch next {
		case end:
			t := token.String{
				Offset: start,
				Mark:   end,
				Value:  result.String(),
			}

			return offset, runes, append(tokens, t), nil

		case '\\':
			if len(runes) == 0 || isNewline(peek(runes, 2)) {
				break
			}

			var (
				code rune
				err  error
			)

			offset, runes, code, err = scanEscape(offset, runes)

			if err != nil {
				return offset, runes, tokens, err
			}

			result.WriteRune(code)

		default:
			result.WriteRune(next)
		}
	}

	t := token.String{
		Offset: start,
		Mark:   end,
		Value:  result.String(),
	}

	return offset, runes, append(tokens, t), SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

// See: https://drafts.csswg.org/css-syntax/#consume-a-name
func scanName(offset int, runes []rune) (int, []rune, string, error) {
	var result strings.Builder

	for len(runes) > 0 {
		switch {
		case isName(peek(runes, 1)):
			result.WriteRune(peek(runes, 1))
			offset, runes = offset+1, runes[1:]

		case startsEscape(runes):
			var (
				code rune
				err  error
			)

			offset, runes, code, err = scanEscape(offset+1, runes[1:])

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
	if isHexDigit(peek(runes, 1)) {
		code := hexValue(peek(runes, 1))
		offset, runes = offset+1, runes[1:]

		for i := 0; len(runes) > 0 && i < 5; i++ {
			if isHexDigit(peek(runes, 1)) {
				code = 0x10*code + hexValue(peek(runes, 1))
				offset, runes = offset+1, runes[1:]
			} else {
				break
			}
		}

		if isWhitespace(peek(runes, 1)) {
			offset, runes = offset+1, runes[1:]
		}

		value := rune(code)

		if isSurrogate(value) || value > 0x10ffff {
			value = 0xfffd
		}

		return offset, runes, value, nil
	}

	value := peek(runes, 1)

	if len(runes) == 0 {
		value = 0xfffd
	}

	return offset + 1, runes[1:], value, nil
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

	if peek(runes, 1) == '(' {
		offset, runes = offset+1, runes[1:]

		if strings.EqualFold(name, "url") {
			for isWhitespace(peek(runes, 1)) {
				offset, runes = offset+1, runes[1:]
			}

			if len(runes) == 0 || (peek(runes, 1) != '"' && peek(runes, 1) != '\'') {
				return scanUrl(offset, runes, tokens, start)
			}
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
func scanUrl(offset int, runes []rune, tokens []token.Token, start int) (int, []rune, []token.Token, error) {
	for isWhitespace(peek(runes, 1)) {
		offset, runes = offset+1, runes[1:]
	}

	var result strings.Builder

	for len(runes) > 0 {
		if isWhitespace(peek(runes, 1)) {
			position := offset

			for isWhitespace(peek(runes, 1)) {
				offset, runes = offset+1, runes[1:]
			}

			t := token.Url{
				Offset: start,
				Value:  result.String(),
			}

			if peek(runes, 1) == ')' {
				offset, runes = offset+1, runes[1:]

				return offset, runes, append(tokens, t), nil
			}

			return offset, runes, tokens, SyntaxError{
				Offset:  position,
				Message: "unexpected whitespace",
			}
		}

		switch peek(runes, 1) {
		case ')':
			t := token.Url{
				Offset: start,
				Value:  result.String(),
			}

			return offset + 1, runes[1:], append(tokens, t), nil

		case '\\':
			if startsEscape(runes) {
				var (
					code rune
					err  error
				)

				offset, runes, code, err = scanEscape(offset+1, runes[1:])

				if err != nil {
					return offset, runes, tokens, err
				}

				result.WriteRune(code)
			} else {
				result.WriteRune(peek(runes, 1))
				offset, runes = offset+1, runes[1:]
			}

		case '"', '\'':
			end := peek(runes, 1)
			offset, runes = offset+1, runes[1:]

			for len(runes) > 0 {
				if peek(runes, 1) == end {
					offset, runes = offset+1, runes[1:]
					break
				}

				result.WriteRune(peek(runes, 1))
				offset, runes = offset+1, runes[1:]
			}

		default:
			result.WriteRune(peek(runes, 1))
			offset, runes = offset+1, runes[1:]
		}
	}

	t := token.Url{
		Offset: start,
		Value:  result.String(),
	}

	return offset, runes, append(tokens, t), SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func isBetween(rune rune, lower rune, upper rune) bool {
	return rune >= lower && rune <= upper
}

// See: https://drafts.csswg.org/css-syntax/#digit
func isDigit(rune rune) bool {
	return isBetween(rune, '0', '9')
}

// See: https://drafts.csswg.org/css-syntax/#hex-digit
func isHexDigit(rune rune) bool {
	return isDigit(rune) || isBetween(rune, 'A', 'F') || isBetween(rune, 'a', 'f')
}

// See: https://drafts.csswg.org/css-syntax/#newline
func isNewline(rune rune) bool {
	switch rune {
	case '\u000A', '\u000C', '\u000D':
		return true
	}

	return false
}

// See: https://drafts.csswg.org/css-syntax/#whitespace
func isWhitespace(rune rune) bool {
	switch rune {
	case '\u0009', '\u0020':
		return true
	}

	return isNewline(rune)
}

// See: https://drafts.csswg.org/css-syntax/#uppercase-letter
func isUppercaseLetter(rune rune) bool {
	return isBetween(rune, 'A', 'Z')
}

// See: https://drafts.csswg.org/css-syntax/#lowercase-letter
func isLowercaseLetter(rune rune) bool {
	return isBetween(rune, 'a', 'z')
}

// See: https://drafts.csswg.org/css-syntax/#letter
func isLetter(rune rune) bool {
	return isUppercaseLetter(rune) || isLowercaseLetter(rune)
}

// See: https://drafts.csswg.org/css-syntax/#non-ascii-code-point
func isAscii(rune rune) bool {
	return rune < '\u0080'
}

// See: https://drafts.csswg.org/css-syntax/#name-start-code-point
func isNameStart(rune rune) bool {
	return isLetter(rune) || !isAscii(rune) || rune == '_'
}

// See: https://drafts.csswg.org/css-syntax/#name-code-point
func isName(rune rune) bool {
	return isNameStart(rune) || isDigit(rune) || rune == '-'
}

// See: https://drafts.csswg.org/css-syntax/#non-printable-code-point
func isNonPrintable(rune rune) bool {
	return isBetween(rune, 0x0000, 0x0008) || rune == 0x000b || isBetween(rune, 0x000e, 0x001f) || rune == 0x007f
}

// See: https://infra.spec.whatwg.org/#surrogate
func isSurrogate(rune rune) bool {
	return isBetween(rune, 0xd800, 0xdfff)
}

func hexValue(rune rune) int {
	if isDigit(rune) {
		return int(rune) - '0'
	}

	if isBetween(rune, 'a', 'f') {
		return int(rune) - 'a' + 10
	}

	if isBetween(rune, 'A', 'F') {
		return int(rune) - 'A' + 10
	}

	return -1
}

// See: https://drafts.csswg.org/css-syntax/#would-start-an-identifier
func startsIdentifier(runes []rune) bool {
	switch peek(runes, 1) {
	case '-':
		return isNameStart(peek(runes, 2)) || peek(runes, 2) == '-' || startsEscape(runes)

	case '\\':
		return startsEscape(runes[1:])
	}

	return isNameStart(peek(runes, 1))
}

// See: https://drafts.csswg.org/css-syntax/#starts-with-a-valid-escape
func startsEscape(runes []rune) bool {
	return peek(runes, 1) == '\\' && !isNewline(peek(runes, 2))
}

// See: https://drafts.csswg.org/css-syntax/#starts-with-a-number
func startsNumber(runes []rune) bool {
	switch peek(runes, 1) {
	case '+', '-':
		return isDigit(peek(runes, 2)) || peek(runes, 2) == '.' && isDigit(peek(runes, 3))
	case '.':
		return isDigit(peek(runes, 2))
	}

	return isDigit(peek(runes, 1))
}

func startsFraction(runes []rune) bool {
	return peek(runes, 1) == '.' && isDigit(peek(runes, 2))
}

func startsExponent(runes []rune) bool {
	if peek(runes, 1) == 'E' || peek(runes, 1) == 'e' {
		if peek(runes, 2) == '+' || peek(runes, 2) == '-' {
			return isDigit(peek(runes, 3))
		}

		return isDigit(peek(runes, 2))
	}

	return false
}

func peek(runes []rune, n int) rune {
	if len(runes) < n {
		return -1
	}

	return runes[n-1]
}
