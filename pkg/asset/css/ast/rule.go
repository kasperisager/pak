package ast

import (
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type Rule interface {
	VisitNode(NodeVisitor)
	VisitRule(RuleVisitor)
}

type RuleVisitor struct {
	AtRule        func(AtRule)
	QualifiedRule func(QualifiedRule)
}

type AtRule struct {
	Name    string
	Prelude []token.Token
	Value   *Block
}

func (r AtRule) VisitNode(v NodeVisitor) {
	v.Rule(r)
}

func (r AtRule) VisitRule(v RuleVisitor) {
	v.AtRule(r)
}

type QualifiedRule struct {
	Prelude []token.Token
	Value   Block
}

func (r QualifiedRule) VisitNode(v NodeVisitor) {
	v.Rule(r)
}

func (r QualifiedRule) VisitRule(v RuleVisitor) {
	v.QualifiedRule(r)
}
