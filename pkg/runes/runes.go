package runes

func IsBetween(rune rune, lower rune, upper rune) bool {
	return rune >= lower && rune <= upper
}

func IsDigit(rune rune) bool {
	return IsBetween(rune, '0', '9')
}

func IsBinaryDigit(rune rune) bool {
	return IsBetween(rune, '0', '1')
}

func IsOctalDigit(rune rune) bool {
	return IsBetween(rune, '0', '7')
}

func IsHexDigit(rune rune) bool {
	return IsDigit(rune) || IsBetween(rune, 'a', 'f') || IsBetween(rune, 'A', 'F')
}

func BinaryValue(rune rune) int {
	if IsBinaryDigit(rune) {
		return int(rune) - '0'
	}

	return -1
}

func OctalValue(rune rune) int {
	if IsOctalDigit(rune) {
		return int(rune) - '0'
	}

	return -1
}

func DecimalValue(rune rune) int {
	if IsDigit(rune) {
		return int(rune) - '0'
	}

	return -1
}

func HexValue(rune rune) int {
	if IsDigit(rune) {
		return int(rune) - '0'
	}

	if IsBetween(rune, 'a', 'f') {
		return int(rune) - 'a' + 10
	}

	if IsBetween(rune, 'A', 'F') {
		return int(rune) - 'A' + 10
	}

	return -1
}
