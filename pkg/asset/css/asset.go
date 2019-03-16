package css

import (
	"bytes"
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/writer"
)

func Asset(url *url.URL, contents []byte) (asset.Asset, error) {
	tokens, err := scanner.Scan(bytes.Runes(contents))

	if err != nil {
		return nil, err
	}

	styleSheet, err := parser.Parse(tokens)

	if err != nil {
		return nil, err
	}

	return &cssAsset{url, styleSheet}, nil
}

type cssAsset struct {
	url *url.URL
	styleSheet ast.StyleSheet
}

func (a *cssAsset) URL() *url.URL {
	return a.url
}

func (a *cssAsset) References() []*url.URL {
	references := []*url.URL{}

	for _, rule := range a.styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			references = append(
				references,
				a.URL().ResolveReference(rule.URL),
			)

		default:
			return references
		}
	}

	return references
}

func (a *cssAsset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.styleSheet)
	return b.Bytes()
}

func (a *cssAsset) Merge(b asset.Asset) bool {
	switch b := b.(type) {
	case *cssAsset:
		needle := b.URL().String()

		for i, rule := range a.styleSheet.Rules {
			switch rule := rule.(type) {
			case ast.ImportRule:
				found := a.URL().ResolveReference(rule.URL).String()

				if needle == found {
					a.styleSheet.Rules = append(
						a.styleSheet.Rules[:i],
						append(
							b.styleSheet.Rules,
							a.styleSheet.Rules[i+1:]...,
						)...,
					)

					return true
				}
			}
		}
	}

	return false
}

func (a *cssAsset) Hoist() {
}
