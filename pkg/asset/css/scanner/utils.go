package scanner

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
	return rune >= 'A' && rune <= 'Z'
}

// See: https://drafts.csswg.org/css-syntax/#lowercase-letter
func isLowercaseLetter(rune rune) bool {
	return rune >= 'a' && rune <= 'z'
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
	n := len(runes)

	if n == 0 {
		return false
	}

	switch runes[0] {
	case '-':
		return n > 1 && isNameStart(runes[1]) || n > 1 && runes[1] == '-' || startsEscape(runes)

	case '\\':
		return startsEscape(runes[1:])
	}

	return isNameStart(runes[0])
}

// See: https://drafts.csswg.org/css-syntax/#starts-with-a-valid-escape
func startsEscape(runes []rune) bool {
	return len(runes) > 1 && runes[0] == '\\' && !isNewline(runes[1])
}

// See: https://drafts.csswg.org/css-syntax/#starts-with-a-number
func startsNumber(runes []rune) bool {
	n := len(runes)

	if n == 0 {
		return false
	}

	switch runes[0] {
	case '+', '-':
		return n > 1 && isDigit(runes[1]) || n > 2 && runes[1] == '.' && isDigit(runes[2])
	case '.':
		return n > 1 && isDigit(runes[1])
	}

	return isDigit(runes[0])
}

func startsFraction(runes []rune) bool {
	return len(runes) > 1 && runes[0] == '.' && isDigit(runes[1])
}

func startsExponent(runes []rune) bool {
	if len(runes) > 1 && (runes[0] == 'E' || runes[0] == 'e') {
		if runes[1] == '+' || runes[1] == '-' {
			return len(runes) > 2 && isDigit(runes[2])
		}

		return isDigit(runes[1])
	}

	return false
}
