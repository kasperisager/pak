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
		importLocation int
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
	base := a.URL()

	var references []asset.Reference

	for _, rule := range a.styleSheet.Rules {
		references = collectReferences(base, rule, references)
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
		n := r.(*cssReference).node

		switch n := n.(type) {
		case ast.ImportRule:
			var ok bool

			a.importLocation, ok = mergeImportRule(
				n,
				&a.styleSheet,
				&b.styleSheet,
				a.importLocation,
			)

			return ok
		}
	}

	return false
}

func (r *cssReference) URL() *url.URL {
	return r.url
}

func collectReferences(
	base *url.URL,
	node ast.Node,
	references []asset.Reference,
) []asset.Reference {
	switch node := node.(type) {
	case ast.StyleSheet:
		for _, rule := range node.Rules {
			references = collectReferences(base, rule, references)
		}

	case ast.StyleRule:
		for _, declaration := range node.Declarations {
			references = collectReferences(base, declaration, references)
		}

	case ast.ImportRule:
		return  append(
			references,
			&cssReference{
				url: base.ResolveReference(node.URL),
				node: node,
			},
		)

	case ast.MediaRule:
		return collectReferences(base, node.StyleSheet, references)

	case ast.SupportsRule:
		return collectReferences(base, node.StyleSheet, references)
	}

	return references
}

func mergeImportRule(
	rule ast.ImportRule,
	source *ast.StyleSheet,
	target *ast.StyleSheet,
	location int,
) (int, bool) {
	for i, r := range source.Rules {
		switch r := r.(type) {
		case ast.ImportRule:
			if r == rule {
				source.Rules = append(
					source.Rules[:i],
					append(
						target.Rules,
						source.Rules[i+1:]...,
					)...,
				)

				return i, true
			}
		}
	}

	return -1, false
}
