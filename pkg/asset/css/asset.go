package css

import (
	"net/url"
	"io"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/writer"
)

func Asset(path string, contents string) (asset.Asset, error) {
	p, err := url.Parse(path)

	if err != nil {
		return nil, err
	}

	tokens, err := scanner.Scan([]rune(contents))

	if err != nil {
		return nil, err
	}

	styleSheet, err := parser.Parse(tokens)

	if err != nil {
		return nil, err
	}

	return cssAsset{*p, styleSheet}, nil
}

type cssAsset struct {
	path       url.URL
	styleSheet ast.StyleSheet
}

func (a cssAsset) Path() string {
	return a.path.String()
}

func (a cssAsset) References() []asset.Reference {
	references := []asset.Reference{}

	for _, rule := range a.styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			reference, err := url.Parse(rule.Url)

			if err != nil {
				return nil
			}

			references = append(references, asset.Reference{
				Path: a.path.ResolveReference(reference).String(),
			})

		default:
			return references
		}
	}

	return references
}

func (a cssAsset) Write(w io.Writer) {
	writer.Write(w, a.styleSheet)
}
