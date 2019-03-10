package css

import (
	"bytes"
	"net/url"
	"path"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/writer"
)

func Asset(filename string, contents []byte) (asset.Asset, error) {
	p, err := url.Parse(filename)

	if err != nil {
		return nil, err
	}

	tokens, err := scanner.Scan(bytes.Runes(contents))

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
	Url        url.URL
	StyleSheet ast.StyleSheet
}

func (a cssAsset) Path() string {
	return a.Url.String()
}

func (a cssAsset) References() []asset.Reference {
	references := []asset.Reference{}

	for _, rule := range a.StyleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			ref, err := url.Parse(rule.Url)

			if err != nil {
				return nil
			}

			if ref.Host != "" {
				continue
			}

			references = append(references, asset.Reference{
				Path: path.Join(path.Dir(a.Url.Path), ref.Path),
			})

		default:
			return references
		}
	}

	return references
}

func (a cssAsset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.StyleSheet)
	return b.Bytes()
}
