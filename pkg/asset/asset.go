package asset

type Asset interface {
	Path() string
	Data() []byte
	References() []Reference
}

type Reference struct {
	Path string
}
