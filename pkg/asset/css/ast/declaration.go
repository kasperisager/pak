package ast

import (
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type Declaration struct {
	Name  string
	Value []token.Token
}
