package writer

import (
	"fmt"
	"io"

	"github.com/kasperisager/pak/pkg/asset/importmap/ast"
)

func Write(w io.Writer, importMap *ast.ImportMap) {
	fmt.Fprintf(w, "{")

	writeImports(w, importMap)
	writeScopes(w, importMap)

	fmt.Fprintf(w, "}")
}

func writeImports(w io.Writer, importMap *ast.ImportMap) {
	if len(importMap.Imports) == 0 {
		return
	}

	fmt.Fprintf(w, `"imports":{`)

	for i, specifier := range importMap.Imports {
		if i != 0 {
			fmt.Fprintf(w, ",")
		}

		writeSpecifier(w, specifier)
	}

	fmt.Fprintf(w, `}`)
}

func writeScopes(w io.Writer, importMap *ast.ImportMap) {
	if len(importMap.Scopes) == 0 {
		return
	}

	fmt.Fprintf(w, `"scopes":{`)

	for i, scope := range importMap.Scopes {
		if i != 0 {
			fmt.Fprintf(w, ",")
		}

		fmt.Fprintf(w, `"%s":{`, scope.Prefix)

		for i, specifier := range scope.Specifiers {
			if i != 0 {
				fmt.Fprintf(w, ",")
			}

			writeSpecifier(w, specifier)
		}

		fmt.Fprintf(w, `}`)
	}

	fmt.Fprintf(w, `}`)
}

func writeSpecifier(w io.Writer, specifier *ast.Specifier) {
	fmt.Fprintf(w, `"%s":[`, specifier.Key)

	for i, address := range specifier.Addresses {
		if i != 0 {
			fmt.Fprintf(w, ",")
		}

		fmt.Fprintf(w, `"%s"`, address)
	}

	fmt.Fprintf(w, `]`)
}
