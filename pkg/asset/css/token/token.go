package token

import (
	"fmt"
)

type Token interface {
	VisitToken(TokenVisitor)
}

type TokenVisitor struct {
	Ident       func(Ident)
	Function    func(Function)
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
	OpenSquare  func(OpenSquare)
	CloseSquare func(CloseSquare)
	OpenParen   func(OpenParen)
	CloseParen  func(CloseParen)
	OpenCurly   func(OpenCurly)
	CloseCurly  func(CloseCurly)
}

type Ident struct {
	Offset int
	Value  string
}

func (t Ident) VisitToken(v TokenVisitor) {
	v.Ident(t)
}

func (t Ident) String() string {
	return t.Value
}

type Function struct {
	Offset int
	Value  string
}

func (t Function) VisitToken(v TokenVisitor) {
	v.Function(t)
}

func (t Function) String() string {
	return t.Value + "("
}

type AtKeyword struct {
	Offset int
	Value  string
}

func (t AtKeyword) VisitToken(v TokenVisitor) {
	v.AtKeyword(t)
}

func (t AtKeyword) String() string {
	return "@" + t.Value
}

type Hash struct {
	Offset int
	Value  string
	Id     bool
}

func (t Hash) VisitToken(v TokenVisitor) {
	v.Hash(t)
}

func (t Hash) String() string {
	return "#" + t.Value
}

type String struct {
	Offset int
	Value  string
}

func (t String) VisitToken(v TokenVisitor) {
	v.String(t)
}

func (t String) String() string {
	return fmt.Sprintf("\"%s\"", t.Value)
}

type Url struct {
	Offset int
	Value  string
}

func (t Url) VisitToken(v TokenVisitor) {
	v.Url(t)
}

func (t Url) String() string {
	return fmt.Sprintf("url(%s)", t.Value)
}

type Delim struct {
	Offset int
	Value  rune
}

func (t Delim) VisitToken(v TokenVisitor) {
	v.Delim(t)
}

func (t Delim) String() string {
	return string(t.Value)
}

type Number struct {
	Offset  int
	Value   float32
	Integer bool
}

func (t Number) VisitToken(v TokenVisitor) {
	v.Number(t)
}

func (t Number) String() string {
	return fmt.Sprintf("%f", t.Value)
}

type Percentage struct {
	Offset int
	Value  float32
}

func (t Percentage) VisitToken(v TokenVisitor) {
	v.Percentage(t)
}

func (t Percentage) String() string {
	return fmt.Sprintf("%f%%", t.Value)
}

type Dimension struct {
	Offset  int
	Value   float32
	Integer bool
	Unit    string
}

func (t Dimension) VisitToken(v TokenVisitor) {
	v.Dimension(t)
}

func (t Dimension) String() string {
	return fmt.Sprintf("%f%s", t.Value, t.Unit)
}

type Whitespace struct {
	Offset int
}

func (t Whitespace) VisitToken(v TokenVisitor) {
	v.Whitespace(t)
}

func (t Whitespace) String() string {
	return "·"
}

type Colon struct {
	Offset int
}

func (t Colon) VisitToken(v TokenVisitor) {
	v.Colon(t)
}

func (t Colon) String() string {
	return ":"
}

type Semicolon struct {
	Offset int
}

func (t Semicolon) VisitToken(v TokenVisitor) {
	v.Semicolon(t)
}

func (t Semicolon) String() string {
	return ";"
}

type Comma struct {
	Offset int
}

func (t Comma) VisitToken(v TokenVisitor) {
	v.Comma(t)
}

func (t Comma) String() string {
	return ","
}

type OpenSquare struct {
	Offset int
}

func (t OpenSquare) VisitToken(v TokenVisitor) {
	v.OpenSquare(t)
}

func (t OpenSquare) String() string {
	return "["
}

type CloseSquare struct {
	Offset int
}

func (t CloseSquare) VisitToken(v TokenVisitor) {
	v.CloseSquare(t)
}

func (t CloseSquare) String() string {
	return "]"
}

type OpenParen struct {
	Offset int
}

func (t OpenParen) VisitToken(v TokenVisitor) {
	v.OpenParen(t)
}

func (t OpenParen) String() string {
	return "("
}

type CloseParen struct {
	Offset int
}

func (t CloseParen) VisitToken(v TokenVisitor) {
	v.CloseParen(t)
}

func (t CloseParen) String() string {
	return ")"
}

type OpenCurly struct {
	Offset int
}

func (t OpenCurly) VisitToken(v TokenVisitor) {
	v.OpenCurly(t)
}

func (t OpenCurly) String() string {
	return "{"
}

type CloseCurly struct {
	Offset int
}

func (t CloseCurly) VisitToken(v TokenVisitor) {
	v.CloseCurly(t)
}

func (t CloseCurly) String() string {
	return "}"
}
