package asset

import (
	"mime"
	"net/url"
	"path"
)

var mediaTypes = map[string]string{
	".css":         "text/css",
	".html":        "text/html",
	".importmap":   "application/importmap+json",
	".js":          "application/javascript",
	".json":        "application/json",
	".png":         "image/png",
	".svg":         "image/svg+xml",
	".webmanifest": "application/manifest+json",
}

func ParseMediaType(value string) string {
	mediaType, _, _ := mime.ParseMediaType(value)
	return mediaType
}

func MediaTypeByExtension(extension string) string {
	return mediaTypes[extension]
}

func MediaTypeByURL(url *url.URL) string {
	return MediaTypeByExtension(path.Ext(url.Path))
}
