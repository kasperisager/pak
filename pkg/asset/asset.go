package asset

import (
	"io"
)

type Asset interface {
	Path() string
	References() []Reference
	Write(io.Writer)
}

type Reference struct {
	Path string
}
