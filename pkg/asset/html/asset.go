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
		url         *url.URL
		flags       asset.Flags
		Attribute   *ast.Attribute
		Conditional bool
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

func (r *Reference) Rewrite(to *url.URL) {
	r.url = to
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
	var flags asset.Flags

	switch element.Name {
	case "link":
		href := element.Attribute("href")

		if href == nil {
			break
		}

		url, err := base.Parse(href.Value)

		if err != nil {
			break
		}

		typ := element.Attribute("type")

		if typ != nil && typ.Value != "" {
			flags = flags.Set("mediaType", typ.Value)
		}

		media := element.Attribute("media")

		conditional := media != nil && media.Value != "" && media.Value != "all"

		rel := element.Attribute("rel")

		if rel == nil {
			break
		}

		switch rel.Value {
		case "stylesheet", "icon", "preload":
			references = append(references, &Reference{
				url:         url,
				flags:       flags,
				Attribute:   href,
				Conditional: conditional,
			})

		case "manifest":
			switch asset.MediaTypeByURL(url) {
			case "application/json":
				flags = flags.Set("mediaType", "application/manifest+json")
			}

			references = append(references, &Reference{
				url:         url,
				flags:       flags,
				Attribute:   href,
				Conditional: conditional,
			})
		}

	case "script":
		src := element.Attribute("src")

		if src == nil {
			break
		}

		url, err := url.Parse(src.Value)

		if err != nil {
			break
		}

		typ := element.Attribute("type")

		if typ == nil {
			break
		}

		switch typ.Value {
		case "importmap":
			references = append(references, &Reference{
				url:   url,
				flags: flags.Set("mediaType", "application/importmap+json"),
				Attribute: src,
			})

		default:
			references = append(references, &Reference{
				url:   url,
				flags: flags.Set("module", typ.Value == "module"),
				Attribute: src,
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
		src := element.Attribute("src")

		if src != nil {
			break
		}

		typ := element.Attribute("type")

		switch {
		case typ != nil && typ.Value == "importmap":
			embeds = append(embeds, &Embed{
				mediaType: "application/importmap+json",
				data:      []byte(element.Text()),
				Element:   element,
			})

		default:
			var flags asset.Flags

			flags = flags.Set("module", typ != nil && typ.Value == "module")

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
