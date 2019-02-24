package ast

type Node interface {
	VisitNode(NodeVisitor)
}

type NodeVisitor struct {
	StyleSheet func(StyleSheet)
	Rule       func(Rule)
}
