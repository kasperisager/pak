package css

import (
	"mime"
	"net/url"
	"path"

	"github.com/kasperisager/pak/pkg/asset"
)

func Asset(url *url.URL, data []byte) (asset.Asset, error) {
	return &BlobAsset{url, data}, nil
}

type BlobAsset struct {
	url  *url.URL
	data []byte
}

func (a *BlobAsset) Type() string {
	return mime.TypeByExtension(path.Ext(a.url.Path))
}

func (a *BlobAsset) URL() *url.URL {
	return a.url
}

func (a *BlobAsset) References() []asset.Reference {
	return nil
}

func (a *BlobAsset) Data() []byte {
	return a.data
}

func (a *BlobAsset) Merge(b asset.Asset, r asset.Reference) bool {
	return false
}
