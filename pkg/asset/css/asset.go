package css

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css/ast"
	"github.com/kasperisager/pak/pkg/asset/css/token"
	"github.com/kasperisager/pak/pkg/asset/css/parser"
	"github.com/kasperisager/pak/pkg/asset/css/scanner"
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

	return &cssAsset{url, styleSheet, 0}, nil
}

type (
	cssAsset struct {
		url        *url.URL
		styleSheet ast.StyleSheet
		importLocation int
	}

	cssImportReference struct {
		url  *url.URL
		rule ast.ImportRule
	}
)

func (a *cssAsset) URL() *url.URL {
	return a.url
}

func (a *cssAsset) References() []asset.Reference {
	return collectReferences(a.URL(), a.styleSheet, nil)
}

func (a *cssAsset) Data() []byte {
	var b bytes.Buffer
	writer.Write(&b, a.styleSheet)
	return b.Bytes()
}

func (a *cssAsset) Merge(b asset.Asset, r asset.Reference) bool {
	switch b := b.(type) {
	case *cssAsset:
		switch r := r.(type) {
		case *cssImportReference:
			var ok bool

			a.importLocation, ok = mergeImportRule(
				r,
				b,
				a,
				a.importLocation,
			)

			return ok
		}
	}

	return false
}

func (r *cssImportReference) URL() *url.URL {
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
				&cssImportReference{
					url: base.ResolveReference(rule.URL),
					rule: rule,
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
	reference *cssImportReference,
	from *cssAsset,
	to *cssAsset,
	location int,
) (int, bool) {
	for i, rule := range to.styleSheet.Rules {
		switch rule := rule.(type) {
		case ast.ImportRule:
			if rule == reference.rule {
				rebaseReferences(&from.styleSheet, from.url, to.url)

				to.styleSheet.Rules = append(
					to.styleSheet.Rules[:i],
					append(
						from.styleSheet.Rules,
						to.styleSheet.Rules[i+1:]...,
					)...,
				)

				return i, true
			}
		}
	}

	return -1, false
}
