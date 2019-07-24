package importmap

import (
	"bytes"
	"net/url"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/importmap/ast"
	"github.com/kasperisager/pak/pkg/asset/importmap/parser"
	"github.com/kasperisager/pak/pkg/asset/importmap/writer"
)

type (
	Asset struct {
		url       *url.URL
		ImportMap *ast.ImportMap
	}

	Reference struct {
		url *url.URL
	}
)

const MediaType = "application/importmap+json"

func From(url *url.URL, data []byte) (*Asset, error) {
	importMap, err := parser.Parse(data)

	if err != nil {
		return nil, err
	}

	return &Asset{
		url:       url,
		ImportMap: importMap,
	}, nil
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
	var b bytes.Buffer
	writer.Write(&b, a.ImportMap)
	return b.Bytes()
}

func (a *Asset) Merge(b asset.Asset, r asset.Relation) bool {
	return false
}

func (r *Reference) URL() *url.URL {
	return r.url
}
