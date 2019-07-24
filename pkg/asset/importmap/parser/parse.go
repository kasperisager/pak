package parser

import (
	"encoding/json"

	"github.com/kasperisager/pak/pkg/asset/importmap/ast"
)

type SyntaxError struct {
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}

func Parse(bytes []byte) (*ast.ImportMap, error) {
	var unmarshalled interface{}

	err := json.Unmarshal(bytes, &unmarshalled)

	if err != nil {
		return nil, err
	}

	entries, ok := unmarshalled.(map[string]interface{})

	if !ok {
		return nil, SyntaxError{
			Message: "import map must be an object",
		}
	}

	importMap := &ast.ImportMap{}

	_, ok = entries["imports"]

	if ok {
		imports, ok := entries["imports"].(map[string]interface{})

		if !ok {
			return nil, SyntaxError{
				Message: `"imports" must be an object`,
			}
		}

		for key, value := range imports {
			specifier := &ast.Specifier{Key: key}

			switch value := value.(type) {
			case string:
				specifier.Addresses = []string{value}

			case []interface{}:
				values := value

				for _, value := range values {
					switch value := value.(type) {
					case string:
						specifier.Addresses = append(specifier.Addresses, value)

					default:
						return nil, SyntaxError{
							Message: `specifier address must be a string`,
						}
					}
				}

			default:
				return nil, SyntaxError{
					Message: `specifier address must be a string or an array of strings`,
				}
			}

			importMap.Imports = append(importMap.Imports, specifier)
		}
	}

	return importMap, nil
}
