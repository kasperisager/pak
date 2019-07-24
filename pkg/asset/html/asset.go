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

type (
	Asset struct {
		url      *url.URL
		Document *ast.Element
	}

	Reference struct {
		url   *url.URL
		flags asset.Flags
	}

	Embed struct {
		mediaType string
		data      []byte
		flags     asset.Flags
		Element   *ast.Element
	}
)

const MediaType = "text/html"

func From(url *url.URL, data []byte) (*Asset, error) {
	runes := bytes.Runes(data)

	tokens, err := scanner.Scan(runes)

	if err != nil {
		return nil, err
	}

	document, err := parser.Parse(tokens)

	if err != nil {
		return nil, err
	}

	return &Asset{
		url:      url,
		Document: document,
	}, nil
}

func (a *Asset) MediaType() string {
	return MediaType
}

func (a *Asset) URL() *url.URL {
	return a.url
}

func (a *Asset) References() []asset.Reference {
	return collectReferences(a.url, a.Document, nil)
}

func (a *Asset) Embeds() []asset.Embed {
	return collectEmbeds(a.url, a.Document, nil)
}

func (a *Asset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.Document)
	return b.Bytes()
}

func (a *Asset) Merge(b asset.Asset, r asset.Relation) bool {
	switch r := r.(type) {
	case *Embed:
		r.Element.Children = []ast.Node{
			&ast.Text{
				Data: string(b.Data()),
			},
		}

		return true
	}

	return false
}

func (r *Reference) VisitRelation(v asset.RelationVisitor) {
	v.Reference(r)
}

func (r *Reference) URL() *url.URL {
	return r.url
}

func (r *Reference) Flags() asset.Flags {
	return r.flags
}

func (e *Embed) VisitRelation(v asset.RelationVisitor) {
	v.Embed(e)
}

func (e *Embed) MediaType() string {
	return e.mediaType
}

func (e *Embed) Data() []byte {
	return e.data
}

func (e *Embed) Flags() asset.Flags {
	return e.flags
}

func collectReferences(
	base *url.URL,
	element *ast.Element,
	references []asset.Reference,
) []asset.Reference {
	switch element.Name {
	case "link":
		rel, _ := element.Attribute("rel")
		href, ok := element.Attribute("href")

		if !ok {
			break
		}

		url, err := url.Parse(href)

		if err != nil {
			break
		}

		switch rel {
		case "stylesheet":
			references = append(references, &Reference{
				url: base.ResolveReference(url),
			})
		}

	case "script":
		typ, _ := element.Attribute("type")
		src, ok := element.Attribute("src")

		if !ok {
			break
		}

		url, err := url.Parse(src)

		if err != nil {
			break
		}

		switch typ {
		case "importmap":
			references = append(references, &Reference{
				url: base.ResolveReference(url),
			})

		default:
			var flags asset.Flags

			flags = flags.Set("module", typ == "module")

			references = append(references, &Reference{
				url:   base.ResolveReference(url),
				flags: flags,
			})
		}
	}

	for _, child := range element.Children {
		switch child := child.(type) {
		case *ast.Element:
			references = collectReferences(base, child, references)
		}
	}

	return references
}

func collectEmbeds(
	base *url.URL,
	element *ast.Element,
	embeds []asset.Embed,
) []asset.Embed {
	switch element.Name {
	case "style":
		embeds = append(embeds, &Embed{
			mediaType: "text/css",
			data:      []byte(element.Text()),
			Element:   element,
		})

	case "script":
		typ, _ := element.Attribute("type")
		_, ok := element.Attribute("src")

		if ok {
			break
		}

		switch typ {
		case "importmap":
			embeds = append(embeds, &Embed{
				mediaType: "application/importmap+json",
				data:      []byte(element.Text()),
				Element:   element,
			})

		default:
			var flags asset.Flags

			flags = flags.Set("module", typ == "module")

			embeds = append(embeds, &Embed{
				mediaType: "application/javascript",
				data:      []byte(element.Text()),
				flags:     flags,
			})
		}
	}

	for _, child := range element.Children {
		switch child := child.(type) {
		case *ast.Element:
			embeds = collectEmbeds(base, child, embeds)
		}
	}

	return embeds
}
