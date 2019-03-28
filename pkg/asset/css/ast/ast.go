package ast

import (
	"net/url"

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
		MediaRule  func(MediaRule)
	}

	StyleRule struct {
		Selectors    []Selector
		Declarations []Declaration
	}

	ImportRule struct {
		URL *url.URL
	}

	MediaRule struct {
		Conditions []MediaQuery
		StyleSheet StyleSheet
	}

	Declaration struct {
		Name      string
		Value     []token.Token
		Important bool
	}

	SelectorMatcher rune

	SelectorModifier int

	SelectorCombinator rune

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

	AttributeSelector struct {
		Name      string
		Namespace *string
		Value     string
		Matcher   SelectorMatcher
		Modifier  SelectorModifier
	}

	TypeSelector struct {
		Name      string
		Namespace *string
	}

	PseudoSelector struct {
		Name       string
		Functional bool
		Value      []token.Token
	}

	CompoundSelector struct {
		Left  Selector
		Right Selector
	}

	RelativeSelector struct {
		Combinator SelectorCombinator
		Left       Selector
		Right      Selector
	}

	MediaQualifier int

	MediaOperator int

	MediaQuery struct {
		Type      string
		Qualifier MediaQualifier
		Condition MediaCondition
	}

	MediaCondition interface {
		VisitMediaCondition(MediaConditionVisitor)
	}

	MediaConditionVisitor struct {
		MediaOperation func(MediaOperation)
		MediaFeature   func(MediaFeature)
		MediaNegation  func(MediaNegation)
	}

	MediaOperation struct {
		Operator MediaOperator
		Left     MediaCondition
		Right    MediaCondition
	}

	MediaFeature struct {
		Name  string
		Value MediaValue
	}

	MediaNegation struct {
		Condition MediaCondition
	}

	MediaValue interface {
		VisitMediaValue(MediaValueVisitor)
	}

	MediaValueVisitor struct {
		MediaValuePlain func(MediaValuePlain)
		MediaValueRange func(MediaValueRange)
	}

	MediaValuePlain struct {
		Value token.Token
	}

	MediaValueRange struct {
		LowerValue     token.Token
		LowerInclusive bool
		UpperValue     token.Token
		UpperInclusive bool
	}
)

const (
	CombinatorDescendant       SelectorCombinator = ' '
	CombinatorDirectDescendant                    = '>'
	CombinatorSibling                             = '~'
	CombinatorDirectSibling                       = '+'
)

const (
	MatcherEqual     SelectorMatcher = '='
	MatcherIncludes                  = '~'
	MatcherDashMatch                 = '|'
	MatcherPrefix                    = '^'
	MatcherSuffix                    = '$'
	MatcherSubstring                 = '*'
)

const (
	QualifierOnly MediaQualifier = iota + 1
	QualifierNot
)

const (
	OperatorAnd MediaOperator = iota + 1
	OperatorOr
)

func (r StyleRule) VisitRule(v RuleVisitor) { v.StyleRule(r) }

func (r ImportRule) VisitRule(v RuleVisitor) { v.ImportRule(r) }

func (r MediaRule) VisitRule(v RuleVisitor) { v.MediaRule(r) }

func (s IdSelector) VisitSelector(v SelectorVisitor) { v.IdSelector(s) }

func (s ClassSelector) VisitSelector(v SelectorVisitor) { v.ClassSelector(s) }

func (s AttributeSelector) VisitSelector(v SelectorVisitor) { v.AttributeSelector(s) }

func (s TypeSelector) VisitSelector(v SelectorVisitor) { v.TypeSelector(s) }

func (s PseudoSelector) VisitSelector(v SelectorVisitor) { v.PseudoSelector(s) }

func (s RelativeSelector) VisitSelector(v SelectorVisitor) { v.RelativeSelector(s) }

func (s CompoundSelector) VisitSelector(v SelectorVisitor) { v.CompoundSelector(s) }

func (s MediaOperation) VisitMediaCondition(v MediaConditionVisitor) { v.MediaOperation(s) }

func (s MediaNegation) VisitMediaCondition(v MediaConditionVisitor) { v.MediaNegation(s) }

func (s MediaFeature) VisitMediaCondition(v MediaConditionVisitor) { v.MediaFeature(s) }

func (s MediaValuePlain) VisitMediaValue(v MediaValueVisitor) { v.MediaValuePlain(s) }
