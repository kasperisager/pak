package asset

import (
	"net/url"
)

type (
	Asset interface {
		URL() *url.URL
		MediaType() string
		Data() []byte
		References() []Reference
		Embeds() []Embed
		Merge(Asset, Relation) bool
	}

	Relation interface {
		VisitRelation(RelationVisitor)
	}

	RelationVisitor struct {
		Reference func(Reference)
		Embed     func(Embed)
	}

	Reference interface {
		Relation
		URL() *url.URL
		Flags() Flags
	}

	Embed interface {
		Relation
		MediaType() string
		Data() []byte
		Flags() Flags
	}
)
