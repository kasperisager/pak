package ast

import (
	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type (
	StyleSheet struct {
		Rules []Rule
	}

	Rule interface {
		VisitRule(RuleVisitor)
	}

	RuleVisitor struct {
		StyleRule  func(StyleRule)
		ImportRule func(ImportRule)
	}

	StyleRule struct {
		Selectors    []Selector
		Declarations []Declaration
	}

	ImportRule struct {
		Url string
	}

	Declaration struct {
		Name      string
		Value     []token.Token
		Important bool
	}

	Selector interface {
		VisitSelector(SelectorVisitor)
	}

	SelectorVisitor struct {
		IdSelector        func(IdSelector)
		ClassSelector     func(ClassSelector)
		AttributeSelector func(AttributeSelector)
		TypeSelector      func(TypeSelector)
		PseudoSelector    func(PseudoSelector)
		CompoundSelector  func(CompoundSelector)
		RelativeSelector  func(RelativeSelector)
	}

	IdSelector struct {
		Name string
	}

	ClassSelector struct {
		Name string
	}

	AttributeMatcher int

	AttributeModifier int

	AttributeSelector struct {
		Name      string
		Namespace *string
		Value     *string
		Matcher   AttributeMatcher
		Modifier  AttributeModifier
	}

	TypeSelector struct {
		Name      string
		Namespace *string
	}

	PseudoSelector struct {
		Name  string
		Value []token.Token
	}

	CompoundSelector struct {
		Left  Selector
		Right Selector
	}

	SelectorCombinator rune

	RelativeSelector struct {
		Combinator SelectorCombinator
		Left       Selector
		Right      Selector
	}
)

const (
	Descendant SelectorCombinator = ' '
	DirectDescendant = '>'
	Sibling = '~'
	DirectSibling = '+'
)

func (r StyleRule) VisitRule(v RuleVisitor) { v.StyleRule(r) }

func (r ImportRule) VisitRule(v RuleVisitor) { v.ImportRule(r) }

func (s IdSelector) VisitSelector(v SelectorVisitor) { v.IdSelector(s) }

func (s ClassSelector) VisitSelector(v SelectorVisitor) { v.ClassSelector(s) }

func (s RelativeSelector) VisitSelector(v SelectorVisitor) { v.RelativeSelector(s) }

func (s CompoundSelector) VisitSelector(v SelectorVisitor) { v.CompoundSelector(s) }
