package asset

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRebase(t *testing.T) {
	assert.Equal(t,
		&url.URL{Path: "../foo/foo.css"},
		rebase(
			&url.URL{Path: "foo.css"},
			&url.URL{Path: "/foo/bar.css"},
			&url.URL{Path: "/bar/baz.css"},
		),
	)

	assert.Equal(t,
		&url.URL{Scheme: "http", Host: "example.com", Path: "/foo/foo.css"},
		rebase(
			&url.URL{Path: "foo.css"},
			&url.URL{Scheme: "http", Host: "example.com", Path: "/foo/bar.css"},
			&url.URL{Path: "/bar/baz.css"},
		),
	)
}

func TestRewrite(t *testing.T) {
	assert.Equal(t,
		&url.URL{Path: "../bar/bar.css"},
		rewrite(
			&url.URL{Path: "/foo/foo.css"},
			&url.URL{Path: "baz.css"},
			&url.URL{Path: "/bar/bar.css"},
		),
	)
}
