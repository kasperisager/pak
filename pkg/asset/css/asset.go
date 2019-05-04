package css

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
	"github.com/kasperisager/pak/pkg/asset/css/token"
	"github.com/kasperisager/pak/pkg/asset/css/writer"
)

func Asset(url *url.URL, data []byte) (asset.Asset, error) {
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

	return &CSSAsset{url, styleSheet}, nil
}

type (
	CSSAsset struct {
		url        *url.URL
		StyleSheet ast.StyleSheet
	}

	CSSImportReference struct {
		url  *url.URL
		Rule ast.ImportRule
	}
)

func (a *CSSAsset) URL() *url.URL {
	return a.url
}

func (a *CSSAsset) References() []asset.Reference {
	return collectReferences(a.URL(), a.StyleSheet, nil)
}

func (a *CSSAsset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.StyleSheet)
	return b.Bytes()
}

func (a *CSSAsset) Merge(b asset.Asset, r asset.Reference) bool {
	switch b := b.(type) {
	case *CSSAsset:
		switch r := r.(type) {
		case *CSSImportReference:
			return mergeImportRule(r, b, a)
		}
	}

	return false
}

func (r *CSSImportReference) URL() *url.URL {
	return r.url
}

func collectReferences(
	base *url.URL,
	styleSheet ast.StyleSheet,
	references []asset.Reference,
) []asset.Reference {
	for _, rule := range styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			return append(
				references,
				&CSSImportReference{
					url:  base.ResolveReference(rule.URL),
					Rule: rule,
				},
			)

		case ast.MediaRule:
			return collectReferences(base, rule.StyleSheet, references)

		case ast.SupportsRule:
			return collectReferences(base, rule.StyleSheet, references)
		}
	}

	return references
}

func rebaseReferences(styleSheet *ast.StyleSheet, from *url.URL, to *url.URL) {
	for _, rule := range styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			*rule.URL = *rebaseUrl(rule.URL, from, to)

		case ast.MediaRule:
			rebaseReferences(&rule.StyleSheet, from, to)
		}
	}
}

func rebaseUrl(reference *url.URL, from *url.URL, to *url.URL) *url.URL {
	if reference.IsAbs() {
		return reference
	}

	from = from.ResolveReference(reference)

	if from.Scheme == to.Scheme && from.Host == to.Host {
		if path.IsAbs(reference.Path) {
			return &url.URL{Path: from.Path}
		}

		path, _ := filepath.Rel(path.Dir(to.Path), from.Path)

		return &url.URL{Path: path}
	}

	return from
}

func mergeImportRule(
	reference *CSSImportReference,
	from *CSSAsset,
	to *CSSAsset,
) bool {
	for i, rule := range to.StyleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			if rule == reference.Rule {
				rebaseReferences(&from.StyleSheet, from.url, to.url)

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
