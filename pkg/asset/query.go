package asset

import (
	"net/url"
)

type Query func(Asset) bool

func ByURL(url *url.URL) Query {
	return func(asset Asset) bool {
		v := asset.URL()
		return v.Scheme == url.Scheme && v.Host == url.Host && v.Path == url.Path
	}
}

func ByMediaType(mediaType string) Query {
	return func(asset Asset) bool {
		return asset.MediaType() == mediaType
	}
}
