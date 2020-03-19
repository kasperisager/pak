package parser

import (
	"github.com/kasperisager/pak/pkg/asset/html/ast"
	"github.com/kasperisager/pak/pkg/asset/html/token"
)

type SyntaxError struct {
	Offset  int
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

func Parse(tokens []token.Token) (*ast.Document, error) {
	offset, tokens, document, err := parseDocument(0, tokens)

	if err != nil {
		return nil, err
	}

	offset, tokens = skipWhitespace(offset, tokens)

	if len(tokens) > 0 {
		return nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token",
		}
	}

	return document, nil
}

func parseDocument(offset int, tokens []token.Token) (int, []token.Token, *ast.Document, error) {
	offset, tokens = skipWhitespace(offset, tokens)

	switch peek(tokens, 1).(type) {
	case token.DocumentType:
		offset, tokens = offset+1, tokens[1:]

	default:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected <!doctype html>",
		}
	}

	offset, tokens, documentElement, err := parseDocumentElement(offset+1, tokens[1:])

	if err != nil {
		return offset, tokens, nil, err
	}

	return offset, tokens, &ast.Document{Root: documentElement}, nil
}

func parseDocumentElement(offset int, tokens []token.Token) (int, []token.Token, *ast.Element, error) {
	documentElement := &ast.Element{Name: "html"}

	offset, tokens = skipWhitespace(offset, tokens)

	switch next := peek(tokens, 1).(type) {
	case token.StartTag:
		switch next.Name {
		case "html":
			documentElement = createElement(next)
			offset, tokens = offset+1, tokens[1:]
		}

	case token.EndTag:
		switch next.Name {
		case "html", "head", "body":
		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected end tag, expected <html> tag",
			}
		}

	case token.DocumentType:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected doctype, expected <html> tag",
		}
	}

	offset, tokens, head, err := parseHead(offset, tokens)

	if err != nil {
		return offset, tokens, nil, err
	}

	offset, tokens, body, err := parseBody(offset, tokens)

	if err != nil {
		return offset, tokens, nil, err
	}

	documentElement.Children = []ast.Node{head, body}

	offset, tokens = skipWhitespace(offset, tokens)

	switch next := peek(tokens, 1).(type) {
	case token.EndTag:
		switch next.Name {
		case "html":
			offset, tokens = offset+1, tokens[1:]
		}
	}

	return offset, tokens, documentElement, nil
}

func parseHead(offset int, tokens []token.Token) (int, []token.Token, *ast.Element, error) {
	head := &ast.Element{Name: "head"}

	offset, tokens = skipWhitespace(offset, tokens)

	switch next := peek(tokens, 1).(type) {
	case token.StartTag:
		switch next.Name {
		case "head":
			head = createElement(next)
			offset, tokens = offset+1, tokens[1:]

		case "html":
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token tag, expected <head> tag",
			}
		}

	case token.EndTag:
		switch next.Name {
		case "html", "head", "body":
		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected <head> tag",
			}
		}

	case token.DocumentType:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected token, expected <html> tag",
		}
	}

	for len(tokens) > 0 {
		var (
			child *ast.Element
			err   error
		)

		offset, tokens, child, err = parseHeadChild(offset, tokens)

		if err != nil {
			return offset, tokens, nil, err
		}

		if child == nil {
			break
		}

		head.Children = append(head.Children, child)
	}

	switch next := peek(tokens, 1).(type) {
	case token.EndTag:
		switch next.Name {
		case "head":
			offset, tokens = offset+1, tokens[1:]
		}
	}

	return offset, tokens, head, nil
}

func parseHeadChild(offset int, tokens []token.Token) (int, []token.Token, *ast.Element, error) {
	offset, tokens = skipWhitespace(offset, tokens)

	switch next := peek(tokens, 1).(type) {
	case token.StartTag:
		switch next.Name {
		case "base", "link", "meta":
			element := createElement(next)
			return offset + 1, tokens[1:], element, nil

		case "title", "script", "style":
			offset, tokens, text, err := parseText(offset+1, tokens[1:])

			if err != nil {
				return offset, tokens, nil, err
			}

			element := createElement(next)

			if text != nil {
				element.Children = append(element.Children, text)
			}

			switch next := peek(tokens, 1).(type) {
			case token.EndTag:
				if next.Name == element.Name {
					return offset + 1, tokens[1:], element, nil
				}

				return offset, tokens, nil, SyntaxError{
					Offset:  offset,
					Message: "unexpected token, expected <" + element.Name + "> tag",
				}
			}
		}
	}

	return offset, tokens, nil, nil
}

func parseBody(offset int, tokens []token.Token) (int, []token.Token, *ast.Element, error) {
	body := &ast.Element{Name: "body"}

	offset, tokens = skipWhitespace(offset, tokens)

	switch next := peek(tokens, 1).(type) {
	case token.StartTag:
		switch next.Name {
		case "body":
			body = createElement(next)
			offset, tokens = offset+1, tokens[1:]

		case "html", "head":
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected <body> tag",
			}
		}

	case token.EndTag:
		switch next.Name {
		case "html", "head", "body":
		default:
			return offset, tokens, nil, SyntaxError{
				Offset:  offset,
				Message: "unexpected token, expected <body> tag",
			}
		}

	case token.DocumentType:
		return offset, tokens, nil, SyntaxError{
			Offset:  offset,
			Message: "unexpected doctype, expected <body> tag",
		}
	}

	for len(tokens) > 0 {
		var (
			child ast.Node
			err   error
		)

		offset, tokens, child, err = parseBodyChild(offset, tokens)

		if err != nil {
			return offset, tokens, nil, err
		}

		if child == nil {
			break
		}

		body.Children = append(body.Children, child)
	}

	switch next := peek(tokens, 1).(type) {
	case token.EndTag:
		switch next.Name {
		case "body":
			offset, tokens = offset+1, tokens[1:]
		}
	}

	return offset, tokens, body, nil
}

func parseBodyChild(offset int, tokens []token.Token) (int, []token.Token, ast.Node, error) {
	switch next := peek(tokens, 1).(type) {
	case token.StartTag:
		element := createElement(next)

		offset, tokens = offset+1, tokens[1:]

		if element.IsVoid() {
			return offset, tokens, element, nil
		}

		for {
			switch next := peek(tokens, 1).(type) {
			case token.EndTag:
				if next.Name == element.Name {
					return offset + 1, tokens[1:], element, nil
				}

				return offset, tokens, nil, SyntaxError{
					Offset:  offset,
					Message: "unexpected token, expected <" + element.Name + "> tag",
				}

			default:
				var (
					child ast.Node
					err   error
				)

				offset, tokens, child, err = parseBodyChild(offset, tokens)

				if err != nil {
					return offset, tokens, nil, err
				}

				if child != nil {
					element.Children = append(element.Children, child)
				}
			}
		}

	case token.Character:
		return parseText(offset, tokens)
	}

	return offset, tokens, nil, nil
}

func parseText(offset int, tokens []token.Token) (int, []token.Token, *ast.Text, error) {
	var runes []rune

	for len(tokens) > 0 {
		switch next := peek(tokens, 1).(type) {
		case token.Character:
			runes = append(runes, next.Data)
			offset, tokens = offset+1, tokens[1:]
			continue
		}

		break
	}

	if runes == nil {
		return offset, tokens, nil, nil
	}

	return offset, tokens, &ast.Text{Data: string(runes)}, nil
}

func createElement(tag token.StartTag) *ast.Element {
	element := &ast.Element{Name: tag.Name}

	for _, attribute := range tag.Attributes {
		element.Attributes = append(element.Attributes, &ast.Attribute{
			Name:  attribute.Name,
			Value: attribute.Value,
		})
	}

	return element
}

func peek(tokens []token.Token, n int) token.Token {
	if len(tokens) < n {
		return nil
	}

	return tokens[n-1]
}

func skipWhitespace(offset int, tokens []token.Token) (int, []token.Token) {
	for {
		switch next := peek(tokens, 1).(type) {
		case token.Character:
			switch next.Data {
			case 0x9, 0xa, 0xc, 0xd, ' ':
				offset, tokens = offset+1, tokens[1:]
				continue
			}
		}

		return offset, tokens
	}
}
