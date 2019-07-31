package asset

import (
	"net/url"
	"path"
	"path/filepath"
)

func rebase(target *url.URL, from *url.URL, to *url.URL) *url.URL {
	if target.IsAbs() {
		return target
	}

	from = from.ResolveReference(target)

	if from.Scheme == to.Scheme && from.Host == to.Host {
		if path.IsAbs(target.Path) {
			return &url.URL{Path: from.Path}
		}

		path, _ := filepath.Rel(
			filepath.FromSlash(path.Dir(to.Path)),
			filepath.FromSlash(from.Path),
		)

		return &url.URL{Path: filepath.ToSlash(path)}
	}

	return from
}

func rewrite(base *url.URL, from *url.URL, to *url.URL) *url.URL {
	if from.Scheme == to.Scheme && from.Host == to.Host {
		if path.IsAbs(from.Path) {
			return &url.URL{Path: to.Path}
		}

		path, _ := filepath.Rel(
			filepath.FromSlash(path.Dir(base.Path)),
			filepath.FromSlash(to.Path),
		)

		return &url.URL{Path: filepath.ToSlash(path)}
	}

	return to
}