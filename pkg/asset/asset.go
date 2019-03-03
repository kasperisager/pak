package asset

type Asset interface {
	Path() string
	References() []Reference
}

type Reference struct {
	Path string
}
