package asset

import (
	"net/url"
)

type Asset interface {
	URL() *url.URL
	References() []*url.URL
	Data() []byte
	Merge(Asset) bool
}
