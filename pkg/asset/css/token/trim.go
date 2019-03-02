package token

func TrimWhitespace(tokens []Token) []Token {
	return TrimLeadingWhitespace(TrimTrailingWhitespace(tokens))
}

func TrimLeadingWhitespace(tokens []Token) []Token {
	if len(tokens) != 0 {
		if _, ok := tokens[0].(Whitespace); ok {
			tokens = tokens[1:]
		}
	}

	return tokens
}

func TrimTrailingWhitespace(tokens []Token) []Token {
	if len(tokens) != 0 {
		if _, ok := tokens[len(tokens)-1].(Whitespace); ok {
			tokens = tokens[:len(tokens)-1]
		}
	}

	return tokens
}
