package html

import (
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
)

func Asset(url *url.URL, data []byte) (asset.Asset, error) {
	return &HTMLAsset{url, data}, nil
}

type HTMLAsset struct {
	url  *url.URL
	data []byte
}

func (a *HTMLAsset) URL() *url.URL {
	return a.url
}

func (a *HTMLAsset) References() []asset.Reference {
	return nil
}

func (a *HTMLAsset) Data() []byte {
	return a.data
}

func (a *HTMLAsset) Merge(b asset.Asset, r asset.Reference) bool {
	return false
}
