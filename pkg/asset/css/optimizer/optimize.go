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
	return atRule
}

func optimizeQualifiedRule(qualifiedRule ast.QualifiedRule) ast.QualifiedRule {
	return qualifiedRule
}
