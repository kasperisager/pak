package optimizer

import (
	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func Optimize(styleSheet ast.StyleSheet) ast.StyleSheet {
	rules := styleSheet.Rules

	for i, rule := range rules {
		rules[i] = optimizeRule(rule)
	}

	return styleSheet
}

func optimizeRule(rule ast.Rule) ast.Rule {
	switch rule := rule.(type) {
	case ast.AtRule:
		return optimizeAtRule(rule)
	case ast.QualifiedRule:
		return optimizeQualifiedRule(rule)
	}

	return rule
}

func optimizeAtRule(atRule ast.AtRule) ast.AtRule {
	atRule.Prelude = trimWhitespace(atRule.Prelude)

	if atRule.Value != nil {
		block := optimizeBlock(*atRule.Value)

		atRule.Value = &block
	}

	return atRule
}

func optimizeQualifiedRule(qualifiedRule ast.QualifiedRule) ast.QualifiedRule {
	qualifiedRule.Prelude = trimWhitespace(qualifiedRule.Prelude)
	qualifiedRule.Value = optimizeBlock(qualifiedRule.Value)

	return qualifiedRule
}

func optimizeBlock(block ast.Block) ast.Block {
	block.Value = trimWhitespace(block.Value)

	return block
}

func trimWhitespace(preserved []ast.Preserved) []ast.Preserved {
	return trimLeadingWhitespace(trimTrailingWhitespace(preserved))
}

func trimLeadingWhitespace(preserved []ast.Preserved) []ast.Preserved {
	if len(preserved) != 0 {
		if _, ok := preserved[0].(ast.Whitespace); ok {
			preserved = preserved[1:]
		}
	}

	return preserved
}

func trimTrailingWhitespace(preserved []ast.Preserved) []ast.Preserved {
	if len(preserved) != 0 {
		if _, ok := preserved[len(preserved)-1].(ast.Whitespace); ok {
			preserved = preserved[:len(preserved)-1]
		}
	}

	return preserved
}
