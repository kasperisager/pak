package html

import (
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
)

func Asset(url *url.URL, data []byte) (asset.Asset, error) {
	return &JSAsset{}, nil
}

type (
	JSAsset struct {
		url *url.URL
	}

	JSReference struct {
		url *url.URL
	}
)

func (a *JSAsset) URL() *url.URL {
	return a.url
}

func (a *JSAsset) References() []asset.Reference {
	return nil
}

func (a *JSAsset) Data() []byte {
	return nil
}

func (a *JSAsset) Merge(b asset.Asset, r asset.Reference) bool {
	return false
}

func (r *JSReference) URL() *url.URL {
	return r.url
}
