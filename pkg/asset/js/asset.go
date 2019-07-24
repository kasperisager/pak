package js

import (
	"bytes"
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/js/ast"
	"github.com/kasperisager/pak/pkg/asset/js/parser"
	"github.com/kasperisager/pak/pkg/asset/js/writer"
)

type (
	Asset struct {
		url     *url.URL
		Program *ast.Program
	}

	Reference struct {
		url         *url.URL
		Declaration *ast.ImportDeclaration
	}
)

const MediaType = "application/javascript"

func From(url *url.URL, data []byte, flags asset.Flags) (*Asset, error) {
	program, err := parser.Parse(bytes.Runes(data))

	if err != nil {
		return nil, err
	}

	return &Asset{url, program}, nil
}

func (a *Asset) MediaType() string {
	return MediaType
}

func (a *Asset) URL() *url.URL {
	return a.url
}

func (a *Asset) References() []asset.Reference {
	return collectReferences(a.url, a.Program, nil)
}

func (a *Asset) Embeds() []asset.Embed {
	return nil
}

func (a *Asset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.Program)
	return b.Bytes()
}

func (a *Asset) Merge(b asset.Asset, r asset.Relation) bool {
	return false
}

func (r *Reference) VisitRelation(v asset.RelationVisitor) {
	v.Reference(r)
}

func (r *Reference) URL() *url.URL {
	return r.url
}

func (r *Reference) Flags() asset.Flags {
	return nil
}

func collectReferences(
	base *url.URL,
	program *ast.Program,
	references []asset.Reference,
) []asset.Reference {
	for _, statement := range program.Body {
		switch statement := statement.(type) {
		case *ast.ImportDeclaration:
			url, err := url.Parse(statement.Source.Value)

			if err == nil {
				references = append(references, &Reference{
					url:         base.ResolveReference(url),
					Declaration: statement,
				})
			}
		}
	}

	return references
}
