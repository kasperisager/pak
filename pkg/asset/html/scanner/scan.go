package scanner

import (
	"strings"

	"github.com/kasperisager/pak/pkg/asset/html/token"
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

// See: https://html.spec.whatwg.org/multipage/parsing.html#data-state
func scanToken(offset int, runes []rune, tokens []token.Token) (int, []rune, []token.Token, error) {
	switch peek(runes, 1) {
	case '<':
		switch peek(runes, 2) {
		case '!':
			if peek(runes, 3) == '-' && peek(runes, 4) == '-' {
				return scanComment(offset+4, runes[4:], tokens, offset)
			}

			return scanDocumentType(offset+2, runes[2:], tokens, offset)

		case '/':
			return scanEndTag(offset+2, runes[2:], tokens, offset)

		default:
			return scanStartTag(offset+1, runes[1:], tokens, offset)
		}
	}

	return offset + 1, runes[1:], append(tokens, token.Character{Offset: offset, Data: peek(runes, 1)}), nil
}

// See: https://html.spec.whatwg.org/multipage/parsing.html#doctype-state
func scanDocumentType(offset int, runes []rune, tokens []token.Token, start int) (int, []rune, []token.Token, error) {
	if len(runes) >= 7 {
		next := strings.ToLower(string(runes[0:7]))

		if next == "doctype" {
			offset, runes = offset+7, runes[7:]

			switch peek(runes, 1) {
			case 0x9, 0xa, 0xc, ' ':
				offset, runes = skipWhitespace(offset+1, runes[1:])

			default:
				return offset, runes, tokens, SyntaxError{
					Offset:  offset,
					Message: `unexpected character, expected space`,
				}
			}

			if len(runes) >= 4 {
				next := string(runes[0:4])

				if next == "html" {
					offset, runes = skipWhitespace(offset+4, runes[4:])

					if peek(runes, 1) == '>' {
						return offset + 1, runes[1:], append(tokens, token.DocumentType{Offset: start}), nil
					}
				}
			}
		} else {
			return offset, runes, tokens, SyntaxError{
				Offset:  offset,
				Message: `unexpected character, expected "doctype"`,
			}
		}
	}

	return offset, runes, tokens, SyntaxError{
		Offset:  offset,
		Message: "unexpected character",
	}
}

// See: https://html.spec.whatwg.org/multipage/parsing.html#comment-start-state
func scanComment(offset int, runes []rune, tokens []token.Token, start int) (int, []rune, []token.Token, error) {
	for len(runes) > 0 {
		if peek(runes, 1) == '-' && peek(runes, 2) == '-' && peek(runes, 3) == '>' {
			return offset + 3, runes[3:], tokens, nil
		}

		offset, runes = offset+1, runes[1:]
	}

	return offset, runes, tokens, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func scanStartTag(offset int, runes []rune, tokens []token.Token, start int) (int, []rune, []token.Token, error) {
	offset, runes, name, err := scanTagName(offset, runes)

	if err != nil {
		return offset, runes, tokens, err
	}

	startTag := token.StartTag{Offset: start, Name: name}

attributes:
	for {
		offset, runes = skipWhitespace(offset, runes)

		switch peek(runes, 1) {
		case '/', '>':
			break attributes

		case '=':
			return offset, runes, tokens, SyntaxError{
				Offset:  offset,
				Message: "unexpected character, expected attribute name",
			}

		default:
			var (
				attribute token.Attribute
				err       error
			)

			offset, runes, attribute, err = scanAttribute(offset, runes)

			if err != nil {
				return offset, runes, tokens, err
			}

			for _, existing := range startTag.Attributes {
				if existing.Name == attribute.Name {
					return offset, runes, tokens, SyntaxError{
						Offset:  attribute.Offset,
						Message: "unexpected duplicate attribute",
					}
				}
			}

			startTag.Attributes = append(startTag.Attributes, attribute)
		}
	}

	switch peek(runes, 1) {
	case '>':
		offset, runes = offset+1, runes[1:]

	case '/':
		offset, runes = offset+1, runes[1:]

		switch peek(runes, 1) {
		case '>':
			offset, runes = offset+1, runes[1:]

		default:
			return offset, runes, tokens, SyntaxError{
				Offset:  offset,
				Message: `unexpected character, expected ">"`,
			}
		}

	default:
		return offset, runes, tokens, SyntaxError{
			Offset:  offset,
			Message: `unexpected character, expected ">"`,
		}
	}

	tokens = append(tokens, startTag)

	switch startTag.Name {
	case "script", "style":
		return scanRawText(offset, runes, tokens, startTag.Name)

	case "title":
		return scanRawData(offset, runes, tokens, startTag.Name)
	}

	return offset, runes, tokens, nil
}

func scanEndTag(offset int, runes []rune, tokens []token.Token, start int) (int, []rune, []token.Token, error) {
	offset, runes, name, err := scanTagName(offset, runes)

	if err != nil {
		return offset, runes, tokens, err
	}

	switch peek(runes, 1) {
	case '>':
		offset, runes = offset+1, runes[1:]

	default:
		return offset, runes, tokens, SyntaxError{
			Offset:  offset,
			Message: `unexpected character, expected ">"`,
		}
	}

	return offset, runes, append(tokens, token.EndTag{Offset: start, Name: name}), nil
}

func scanTagName(offset int, runes []rune) (int, []rune, string, error) {
	end := 0

	for len(runes) > end {
		next := peek(runes, end+1)

		switch next {
		case 0x9, 0xa, 0xc, ' ', '/', '>':
			return offset + end, runes[end:], string(runes[0:end]), nil

		case 0x0:
			return offset + end, runes[end:], "", SyntaxError{
				Offset:  offset,
				Message: "unexpected null character",
			}
		}

		end++
	}

	return offset + end, runes[end:], "", SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func scanAttribute(offset int, runes []rune) (int, []rune, token.Attribute, error) {
	attribute := token.Attribute{Offset: offset}

	var (
		name string
		err  error
	)

	offset, runes, name, err = scanAttributeName(offset, runes)

	if err != nil {
		return offset, runes, attribute, err
	}

	attribute.Name = name

	offset, runes = skipWhitespace(offset, runes)

	if peek(runes, 1) == '=' {
		var value string

		offset, runes, value, err = scanAttributeValue(offset+1, runes[1:])

		if err != nil {
			return offset, runes, attribute, err
		}

		attribute.Value = value
	}

	return offset, runes, attribute, nil
}

// See: https://html.spec.whatwg.org/multipage/parsing.html#attribute-name-state
func scanAttributeName(offset int, runes []rune) (int, []rune, string, error) {
	end := 0

	for len(runes) > end {
		next := peek(runes, end+1)

		switch next {
		case 0x9, 0xa, 0xc, ' ', '/', '>', '=':
			return offset + end, runes[end:], strings.ToLower(string(runes[0:end])), nil

		case 0x0:
			return offset + end, runes[end:], "", SyntaxError{
				Offset:  offset,
				Message: "unexpected null character",
			}

		case '"', '\'', '<':
			return offset + end, runes[end:], "", SyntaxError{
				Offset:  offset,
				Message: "unexpected character",
			}

		default:
			end++
		}
	}

	return offset + end, runes[end:], "", SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

// See: https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-value-state
func scanAttributeValue(offset int, runes []rune) (int, []rune, string, error) {
	offset, runes = skipWhitespace(offset, runes)

	switch next := peek(runes, 1); next {
	case '"', '\'':
		offset, runes = offset+1, runes[1:]

		mark := next
		end := 0

		for len(runes) > 0 {
			switch peek(runes, end+1) {
			case mark:
				return offset + end + 1, runes[end+1:], string(runes[0:end]), nil

			case 0x0:
				return offset + end, runes[end:], "", SyntaxError{
					Offset:  offset,
					Message: "unexpected null character",
				}

			default:
				end++
			}
		}

	case '>':
		return offset, runes, "", SyntaxError{
			Offset:  offset,
			Message: "unexpected character, expected attribute value",
		}

	default:
		end := 0

		for len(runes) > 0 {
			switch peek(runes, end+1) {
			case 0x9, 0xa, 0xc, ' ', '>':
				return offset + end, runes[end:], string(runes[0:end]), nil

			case 0x0:
				return offset + end, runes[end:], "", SyntaxError{
					Offset:  offset,
					Message: "unexpected null character",
				}

			case '"', '\'', '<', '=', '`':
				return offset + end, runes[end:], "", SyntaxError{
					Offset:  offset,
					Message: "unexpected character",
				}

			default:
				end++
			}
		}

	}

	return offset, runes, "", SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

// See: https://html.spec.whatwg.org/multipage/parsing.html#rawtext-state
func scanRawText(offset int, runes []rune, tokens []token.Token, tagName string) (int, []rune, []token.Token, error) {
	for len(runes) > 0 {
		if peek(runes, 1) == '<' && peek(runes, 2) == '/' && isLetter(peek(runes, 3)) {
			break
		}

		tokens = append(tokens, token.Character{Offset: offset, Data: runes[0]})
		offset, runes = offset+1, runes[1:]
	}

	end := 2

	for len(runes) > end {
		next := peek(runes, end+1)

		switch next {
		case 0x9, 0xa, 0xc, ' ', '/', '>':
			if tagName == strings.ToLower(string(runes[2:end])) {
				return offset, runes, tokens, nil
			}

		default:
			if !isLetter(next) {
				for i := 0; i < end; i++ {
					tokens = append(tokens, token.Character{Offset: offset, Data: runes[0]})
					offset, runes = offset+1, runes[1:]
				}

				return scanRawText(offset, runes, tokens, tagName)
			}

			end++
		}
	}

	return offset + end, runes[end:], tokens, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

// See: https://html.spec.whatwg.org/multipage/parsing.html#rcdata-state
func scanRawData(offset int, runes []rune, tokens []token.Token, tagName string) (int, []rune, []token.Token, error) {
	for len(runes) > 0 {
		if peek(runes, 1) == '<' && peek(runes, 2) == '/' && isLetter(peek(runes, 3)) {
			break
		}

		tokens = append(tokens, token.Character{Offset: offset, Data: runes[0]})
		offset, runes = offset+1, runes[1:]
	}

	end := 2

	for len(runes) > end {
		next := peek(runes, end+1)

		switch next {
		case 0x9, 0xa, 0xc, ' ', '/', '>':
			if tagName == strings.ToLower(string(runes[2:end])) {
				return offset, runes, tokens, nil
			}

		default:
			if !isLetter(next) {
				for i := 0; i < end; i++ {
					tokens = append(tokens, token.Character{Offset: offset, Data: runes[0]})
					offset, runes = offset+1, runes[1:]
				}

				return scanRawText(offset, runes, tokens, tagName)
			}

			end++
		}
	}

	return offset + end, runes[end:], tokens, SyntaxError{
		Offset:  offset,
		Message: "unexpected end of file",
	}
}

func isBetween(rune rune, lower rune, upper rune) bool {
	return rune >= lower && rune <= upper
}

// See: https://infra.spec.whatwg.org/#ascii-upper-alpha
func isUppercaseLetter(rune rune) bool {
	return isBetween(rune, 'A', 'Z')
}

// See: https://infra.spec.whatwg.org/#ascii-lower-alpha
func isLowercaseLetter(rune rune) bool {
	return isBetween(rune, 'a', 'z')
}

// See: https://infra.spec.whatwg.org/#ascii-alpha
func isLetter(rune rune) bool {
	return isUppercaseLetter(rune) || isLowercaseLetter(rune)
}

func skipWhitespace(offset int, runes []rune) (int, []rune) {
	for {
		switch peek(runes, 1) {
		case 0x9, 0xa, 0xc, ' ':
			offset, runes = offset+1, runes[1:]

		default:
			return offset, runes
		}
	}
}

func peek(runes []rune, n int) rune {
	if len(runes) < n {
		return -1
	}

	return runes[n-1]
}
