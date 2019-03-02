package parser

import (
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type SyntaxError struct {
	Token   token.Token
	Message string
}

func (err SyntaxError) Error() string {
	return err.Message
}
