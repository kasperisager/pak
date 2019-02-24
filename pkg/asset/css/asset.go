package css

import (
	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/optimizer"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
)

func Asset(path string, contents string) (asset.Asset, error) {
	tokens, err := scanner.Scan([]rune(contents))

	if err != nil {
		return nil, err
	}

	styleSheet, err := parser.Parse(tokens)

	if err != nil {
		return nil, err
	}

	asset := cssAsset{
		path,
		optimizer.Optimize(styleSheet),
	}

	return asset, nil
}

type cssAsset struct {
	path       string
	styleSheet ast.StyleSheet
}

func (a cssAsset) Path() string {
	return a.path
}

func (a cssAsset) Rename(path string) {
}

func (a cssAsset) References() []asset.Reference {
	references := []asset.Reference{}

	for _, rule := range a.styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.AtRule:
			switch rule.Name {
			case "import":
				if len(rule.Prelude) > 0 {
					switch value := rule.Prelude[0].(type) {
					case ast.String:
						references = append(references, asset.Reference{Path: value.Value})
					case ast.Url:
						references = append(references, asset.Reference{Path: value.Value})
					}
				}

			case "charset":
				continue

			default:
				return references
			}

		default:
			return references
		}
	}

	return references
}
