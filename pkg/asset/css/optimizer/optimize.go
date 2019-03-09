package optimizer

import (
	"github.com/kasperisager/pak/pkg/asset/css/ast"
)

func Optimize(styleSheet ast.StyleSheet) ast.StyleSheet {
	rules := make([]ast.Rule, len(styleSheet.Rules))

	for i, n := 0, len(rules); i < n; i++ {
		rules[i] = optimizeRule(styleSheet.Rules[i])
	}

	styleSheet.Rules = rules

	return styleSheet
}

func optimizeRule(rule ast.Rule) ast.Rule {
	return rule
}
