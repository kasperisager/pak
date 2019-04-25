package css

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/token"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/writer"
)

func Asset(url *url.URL, data []byte) (asset.Asset, error) {
	runes := bytes.Runes(data)

	tokens, err := scanner.Scan(runes)

	if err != nil {
		return nil, err
	}

	styleSheet, err := parser.Parse(tokens)

	if err != nil {
		switch err := err.(type) {
		case parser.SyntaxError:
			t := tokens[err.Offset]
			return nil, fmt.Errorf("%s: %#v %d", err, t, token.Offset(t))
		}

		return nil, err
	}

	return &cssAsset{url, styleSheet, 0}, nil
}

type (
	cssAsset struct {
		url        *url.URL
		styleSheet ast.StyleSheet
		hoistIndex int
	}

	cssReference struct {
		url  *url.URL
		node ast.Node
	}
)

func (a *cssAsset) URL() *url.URL {
	return a.url
}

func (a *cssAsset) References() []asset.Reference {
	var references []asset.Reference

	for _, rule := range a.styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			references = append(
				references,
				&cssReference{
					url: a.URL().ResolveReference(rule.URL),
					node: rule,
				},
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

func (a *cssAsset) Merge(b asset.Asset, r asset.Reference) bool {
	switch b := b.(type) {
	case *cssAsset:
		for i, rule := range a.styleSheet.Rules {
			switch rule := rule.(type) {
			case ast.ImportRule:
				if r.(*cssReference).node == rule {
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

func (r *cssReference) URL() *url.URL {
	return r.url
}
