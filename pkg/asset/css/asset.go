package css

import (
	"bytes"
	"fmt"
	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/token"
	"github.com/kasperisager/pak/pkg/asset/css/writer"
	"net/url"
)

type (
	Asset struct {
		url        *url.URL
		StyleSheet *ast.StyleSheet
	}

	Reference struct {
		url  *url.URL
		Rule ast.Rule
	}
)

const MediaType = "text/css"

func From(url *url.URL, data []byte) (*Asset, error) {
	runes := bytes.Runes(data)

	tokens, err := scanner.Scan(runes)

	if err != nil {
		return nil, err
	}

	styleSheet, err := parser.Parse(tokens)

	if err != nil {
		switch err := err.(type) {
		case parser.SyntaxError:
			t := tokens[err.Offset]
			return nil, fmt.Errorf("%s: %#v %d", err, t, token.Offset(t))
		}

		return nil, err
	}

	return &Asset{url, styleSheet}, nil
}

func (a *Asset) MediaType() string {
	return MediaType
}

func (a *Asset) URL() *url.URL {
	return a.url
}

func (a *Asset) References() []asset.Reference {
	return collectReferences(a.url, a.StyleSheet, nil)
}

func (a *Asset) Embeds() []asset.Embed {
	return nil
}

func (a *Asset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.StyleSheet)
	return b.Bytes()
}

func (a *Asset) Merge(b asset.Asset, r asset.Relation) bool {
	switch b := b.(type) {
	case *Asset:
		switch r := r.(type) {
		case *Reference:
			return mergeRule(r.Rule, b, a)
		}
	}

	return false
}

func (r *Reference) VisitRelation(v asset.RelationVisitor) {
	v.Reference(r)
}

func (r *Reference) URL() *url.URL {
	return r.url
}

func (r *Reference) Rewrite(to *url.URL) {
	r.url = to

	switch rule := r.Rule.(type) {
	case *ast.ImportRule:
		rule.URL = to
	}
}

func (r *Reference) Flags() asset.Flags {
	return nil
}

func collectReferences(
	base *url.URL,
	styleSheet *ast.StyleSheet,
	references []asset.Reference,
) []asset.Reference {
	for _, rule := range styleSheet.Rules {
		switch rule := rule.(type) {
		case *ast.ImportRule:
			references = append(
				references,
				&Reference{
					url:  rule.URL,
					Rule: rule,
				},
			)

		case *ast.MediaRule:
			collectReferences(base, rule.StyleSheet, references)

		case *ast.SupportsRule:
			collectReferences(base, rule.StyleSheet, references)
		}
	}

	return references
}

func mergeRule(rule ast.Rule, from *Asset, to *Asset) bool {
	for i, found := range to.StyleSheet.Rules {
		if found == rule {
			switch rule.(type) {
			case *ast.ImportRule:
				to.StyleSheet.Rules = append(
					to.StyleSheet.Rules[:i],
					append(
						from.StyleSheet.Rules,
						to.StyleSheet.Rules[i+1:]...,
					)...,
				)

				return true
			}
		}
	}

	return false
}
