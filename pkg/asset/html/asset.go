package html

import (
	"bytes"
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/html/ast"
	"github.com/kasperisager/pak/pkg/asset/html/parser"
	"github.com/kasperisager/pak/pkg/asset/html/scanner"
	"github.com/kasperisager/pak/pkg/asset/html/writer"
)

func Asset(url *url.URL, data []byte) (asset.Asset, error) {
	runes := bytes.Runes(data)

	tokens, err := scanner.Scan(runes)

	if err != nil {
		return nil, err
	}

	document, err := parser.Parse(tokens)

	if err != nil {
		return nil, err
	}

	return &HTMLAsset{url, document}, nil
}

type (
	HTMLAsset struct {
		url      *url.URL
		Document ast.Element
	}

	HTMLReference struct {
		url *url.URL
	}
)

func (a *HTMLAsset) URL() *url.URL {
	return a.url
}

func (a *HTMLAsset) References() []asset.Reference {
	return collectReferences(a.url, a.Document, nil)
}

func (a *HTMLAsset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.Document)
	return b.Bytes()
}

func (a *HTMLAsset) Merge(b asset.Asset, r asset.Reference) bool {
	return false
}

func (r *HTMLReference) URL() *url.URL {
	return r.url
}

func collectReferences(
	base *url.URL,
	element ast.Element,
	references []asset.Reference,
) []asset.Reference {
	switch element.Name {
	case "link":
		rel, _ := element.Attribute("rel")
		href, _ := element.Attribute("href")

		if rel == "stylesheet" {
			href, err := url.Parse(href)

			if err == nil {
				references = append(
					references,
					&HTMLReference{url: base.ResolveReference(href)},
				)
			}
		}
	}

	for _, child := range element.Children {
		switch child := child.(type) {
		case ast.Element:
			references = collectReferences(base, child, references)
		}
	}

	return references
}
