package ast

type (
	Specifier struct {
		Key       string
		Addresses []string
	}

	Scope struct {
		Prefix     string
		Specifiers []*Specifier
	}

	ImportMap struct {
		Imports []*Specifier
		Scopes  []*Scope
	}
)
