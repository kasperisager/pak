package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/html/ast"
)

func Write(w io.Writer, document ast.Element) {
	fmt.Fprintf(w, "<!doctype html>")
	writeElement(w, document)
}

func writeElement(w io.Writer, element ast.Element) {
	fmt.Fprintf(w, "<%s", element.Name)

	for _, attribute := range element.Attributes {
		fmt.Fprintf(w, " ")

		writeAttribute(w, attribute)
	}

	fmt.Fprintf(w, ">")

	for _, child := range element.Children {
		switch child := child.(type) {
		case ast.Element:
			writeElement(w, child)

		case ast.Text:
			writeText(w, child)
		}
	}

	if !element.IsVoid() {
		fmt.Fprintf(w, "</%s>", element.Name)
	}
}

func writeAttribute(w io.Writer, attribute ast.Attribute) {
	fmt.Fprintf(w, attribute.Name)

	if attribute.Value != "" {
		fmt.Fprintf(w, `="%s"`, attribute.Value)
	}
}

func writeText(w io.Writer, text ast.Text) {
	fmt.Fprintf(w, text.Data)
}
