package webmanifest

import (
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/webmanifest/ast"
)

type (
	Asset struct {
		url         *url.URL
		WebManifest *ast.WebManifest
	}

	Reference struct {
		url *url.URL
	}
)

const MediaType = "application/webmanifest+json"

func From(url *url.URL, data []byte, flags asset.Flags) (*Asset, error) {
	return &Asset{url: url}, nil
}

func (a *Asset) MediaType() string {
	return MediaType
}

func (a *Asset) URL() *url.URL {
	return a.url
}

func (a *Asset) References() []asset.Reference {
	return nil
}

func (a *Asset) Embeds() []asset.Embed {
	return nil
}

func (a *Asset) Data() []byte {
	return nil
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

func (r *Reference) Rewrite(to *url.URL) {
}

func (r *Reference) Flags() asset.Flags {
	return nil
}
