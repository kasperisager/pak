package asset

import (
	"net/url"
)

type Asset interface {
	URL() *url.URL
	References() []Reference
	Data() []byte
	Merge(Asset, Reference) bool
}

type Reference interface {
	URL() *url.URL
}
