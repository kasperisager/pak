package css

import (
	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
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

	return cssAsset{path, styleSheet}, nil
}

type cssAsset struct {
	path       string
	styleSheet ast.StyleSheet
}

func (a cssAsset) Path() string {
	return a.path
}

func (a cssAsset) References() []asset.Reference {
	references := []asset.Reference{}

	for _, rule := range a.styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			references = append(references, asset.Reference{Path: rule.Url})

		default:
			return references
		}
	}

	return references
}
