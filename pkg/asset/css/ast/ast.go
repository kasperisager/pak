package ast

import (
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type (
	Node interface {
		VisitNode(NodeVisitor)
	}

	NodeVisitor struct {
		StyleSheet func(StyleSheet)
		Rule       func(Rule)
	}

	StyleSheet struct {
		Rules []Rule
	}

	Rule interface {
		Node
		VisitRule(RuleVisitor)
	}

	RuleVisitor struct {
		AtRule        func(AtRule)
		QualifiedRule func(QualifiedRule)
	}

	AtRule struct {
		Name    string
		Prelude []token.Token
		Value   *Block
	}

	QualifiedRule struct {
		Prelude []token.Token
		Value   Block
	}

	Block struct {
		Value []token.Token
	}

	Declaration struct {
		Name  string
		Value []token.Token
	}
)

func (s StyleSheet) VisitNode(v NodeVisitor) { v.StyleSheet(s) }

func (r AtRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r AtRule) VisitRule(v RuleVisitor) { v.AtRule(r) }

func (r QualifiedRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r QualifiedRule) VisitRule(v RuleVisitor) { v.QualifiedRule(r) }
