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
		AtRule        func(AtRule)
		QualifiedRule func(QualifiedRule)
	}

	AtRule struct {
		Name    string
		Prelude []token.Token
		Value   *Block
	}

	ImportRule struct {
		AtRule
		Url string
	}

	QualifiedRule struct {
		Prelude []token.Token
		Value   Block
	}

	StyleRule struct {
		QualifiedRule
		Selector     Selector
		Declarations []Declaration
	}

	Block struct {
		Value []token.Token
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

	SelectorCombinator int

	RelativeSelector struct {
		Combinator SelectorCombinator
		Left       Selector
		Right      Selector
	}
)

func (r AtRule) VisitRule(v RuleVisitor) { v.AtRule(r) }

func (r QualifiedRule) VisitRule(v RuleVisitor) { v.QualifiedRule(r) }
