package writer

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kasperisager/pak/pkg/asset/js/ast"
)

func Write(w io.Writer, program *ast.Program) {
	writeProgram(w, program)
}

func writeProgram(w io.Writer, program *ast.Program) {
	for i, statement := range program.Body {
		if i != 0 {
			fmt.Fprintf(w, ";")
		}

		switch statement := statement.(type) {
		case ast.Statement:
			writeStatement(w, statement)

		case ast.ModuleDeclaration:
			writeModuleDeclaration(w, statement)
		}
	}
}

func writeLiteral(w io.Writer, literal ast.Literal) {
	switch literal := literal.(type) {
	case *ast.StringLiteral:
		fmt.Fprintf(w, "%q", literal.Value)

	case *ast.BooleanLiteral:
		fmt.Fprintf(w, "%t", literal.Value)

	case *ast.NullLiteral:
		fmt.Fprintf(w, "null")

	case *ast.NumberLiteral:
		fmt.Fprintf(w, "%s", strconv.FormatFloat(literal.Value, 'f', -1, 64))

	case *ast.RegExpLiteral:
		fmt.Fprintf(w, "/%s/%s", literal.Regex.Pattern, literal.Regex.Flags)
	}
}

func writeIdentifier(w io.Writer, identifier *ast.Identifier) {
	fmt.Fprintf(w, "%s", identifier.Name)
}

func writeStatement(w io.Writer, statement ast.Statement) {
	switch statement := statement.(type) {
	case *ast.ExpressionStatement:
		writeExpression(w, statement.Expression)
	}
}

func writeExpression(w io.Writer, expression ast.Expression) {
	switch expression := expression.(type) {
	case ast.Literal:
		writeLiteral(w, expression)
	}
}

func writeModuleDeclaration(w io.Writer, moduleDeclaration ast.ModuleDeclaration) {
	switch moduleDeclaration := moduleDeclaration.(type) {
	case *ast.ImportDeclaration:
		fmt.Fprintf(w, "import")

		specifiers := moduleDeclaration.Specifiers

		if len(specifiers) > 0 {
			defaultImport, ok := specifiers[0].(*ast.ImportDefaultSpecifier)

			if ok {
				fmt.Fprintf(w, " ")
				writeIdentifier(w, defaultImport.Local)

				specifiers = specifiers[1:]

				if len(specifiers) > 0 {
					fmt.Fprintf(w, ",")
				}
			}
		}

		if len(specifiers) > 0 {
			fmt.Fprintf(w, "{")

			for i, specifier := range specifiers {
				if i != 0 {
					fmt.Fprintf(w, ",")
				}

				switch specifier := specifier.(type) {
				case *ast.ImportSpecifier:
					local, imported := specifier.Local, specifier.Imported

					writeIdentifier(w, local)

					if local.Name != imported.Name {
						fmt.Fprintf(w, " as ")
						writeIdentifier(w, imported)
					}
				}
			}

			fmt.Fprintf(w, "}")
		} else {
			fmt.Fprintf(w, " ")
		}

		fmt.Fprintf(w, "from")

		writeLiteral(w, moduleDeclaration.Source)
	}
}
