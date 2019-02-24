package ast

type Rule interface {
	VisitRule(RuleVisitor)
}

type RuleVisitor struct {
	AtRule        func(AtRule)
	QualifiedRule func(QualifiedRule)
}

type AtRule struct {
	Name    string
	Prelude []Preserved
	Value   *Block
}

func (r AtRule) VisitNode(v NodeVisitor) {
	v.Rule(r)
}

func (r AtRule) VisitRule(v RuleVisitor) {
	v.AtRule(r)
}

type QualifiedRule struct {
	Prelude []Preserved
	Value   Block
}

func (r QualifiedRule) VisitNode(v NodeVisitor) {
	v.Rule(r)
}

func (r QualifiedRule) VisitRule(v RuleVisitor) {
	v.QualifiedRule(r)
}
