package ast

import ()

type StyleSheet struct {
	Rules []Rule
}

func (s StyleSheet) VisitNode(v NodeVisitor) {
	v.StyleSheet(s)
}
