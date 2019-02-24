package ast

import (
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type Preserved interface {
	VisitPreserved(PreservedVisitor)
}

type PreservedVisitor struct {
	Ident       func(Ident)
	AtKeyword   func(AtKeyword)
	Hash        func(Hash)
	String      func(String)
	Url         func(Url)
	Delim       func(Delim)
	Number      func(Number)
	Percentage  func(Percentage)
	Dimension   func(Dimension)
	Whitespace  func(Whitespace)
	Colon       func(Colon)
	Semicolon   func(Semicolon)
	Comma       func(Comma)
	CloseSquare func(CloseSquare)
	CloseParen  func(CloseParen)
	CloseCurly  func(CloseCurly)
}

type Ident struct{ token.Ident }

func (t Ident) VisitPreserved(v PreservedVisitor) {
	v.Ident(t)
}

type AtKeyword struct{ token.AtKeyword }

func (t AtKeyword) VisitPreserved(v PreservedVisitor) {
	v.AtKeyword(t)
}

type Hash struct{ token.Hash }

func (t Hash) VisitPreserved(v PreservedVisitor) {
	v.Hash(t)
}

type String struct{ token.String }

func (t String) VisitPreserved(v PreservedVisitor) {
	v.String(t)
}

type Url struct{ token.Url }

func (t Url) VisitPreserved(v PreservedVisitor) {
	v.Url(t)
}

type Delim struct{ token.Delim }

func (t Delim) VisitPreserved(v PreservedVisitor) {
	v.Delim(t)
}

type Number struct{ token.Number }

func (t Number) VisitPreserved(v PreservedVisitor) {
	v.Number(t)
}

type Percentage struct{ token.Percentage }

func (t Percentage) VisitPreserved(v PreservedVisitor) {
	v.Percentage(t)
}

type Dimension struct{ token.Dimension }

func (t Dimension) VisitPreserved(v PreservedVisitor) {
	v.Dimension(t)
}

type Whitespace struct{ token.Whitespace }

func (t Whitespace) VisitPreserved(v PreservedVisitor) {
	v.Whitespace(t)
}

type Colon struct{ token.Colon }

func (t Colon) VisitPreserved(v PreservedVisitor) {
	v.Colon(t)
}

type Semicolon struct{ token.Semicolon }

func (t Semicolon) VisitPreserved(v PreservedVisitor) {
	v.Semicolon(t)
}

type Comma struct{ token.Comma }

func (t Comma) VisitPreserved(v PreservedVisitor) {
	v.Comma(t)
}

type CloseSquare struct{ token.CloseSquare }

func (t CloseSquare) VisitPreserved(v PreservedVisitor) {
	v.CloseSquare(t)
}

type CloseParen struct{ token.CloseParen }

func (t CloseParen) VisitPreserved(v PreservedVisitor) {
	v.CloseParen(t)
}

type CloseCurly struct{ token.CloseCurly }

func (t CloseCurly) VisitPreserved(v PreservedVisitor) {
	v.CloseCurly(t)
}
