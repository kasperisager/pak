package blob

import (
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
)

type Asset struct {
	url  *url.URL
	data []byte
}

func From(url *url.URL, data []byte) *Asset {
	return &Asset{url, data}
}

func (a *Asset) MediaType() string {
	return asset.MediaTypeByURL(a.url)
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
	return a.data
}

func (a *Asset) Merge(b asset.Asset, r asset.Relation) bool {
	return false
}
