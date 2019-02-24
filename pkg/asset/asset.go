package asset

type Asset interface {
	Path() string
	Rename(path string)
	References() []Reference
}

type Reference struct {
	Path string
}
