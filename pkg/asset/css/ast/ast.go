package ast

import (
	"net/url"

	"github.com/kasperisager/pak/pkg/asset/css/token"
)

type (
	Node interface {
		VisitNode(NodeVisitor)
	}

	NodeVisitor struct {
		StyleSheet  func(StyleSheet)
		Rule        func(Rule)
		Declaration func(Declaration)
		Selector    func(Selector)
	}

	StyleSheet struct {
		Node
		Rules []Rule
	}

	Rule interface {
		Node
		VisitRule(RuleVisitor)
	}

	RuleVisitor struct {
		StyleRule     func(StyleRule)
		ImportRule    func(ImportRule)
		MediaRule     func(MediaRule)
		KeyframesRule func(KeyframesRule)
		SupportsRule  func(SupportsRule)
		PageRule      func(PageRule)
	}

	StyleRule struct {
		Rule
		Selectors    []Selector
		Declarations []Declaration
	}

	ImportRule struct {
		Rule
		URL *url.URL
	}

	MediaRule struct {
		Rule
		Conditions []MediaQuery
		StyleSheet StyleSheet
	}

	KeyframesRule struct {
		Rule
		Prefix string
		Name   string
		Blocks []KeyframeBlock
	}

	SupportsRule struct {
		Rule
		Condition  SupportsCondition
		StyleSheet StyleSheet
	}

	PageRule struct {
		Rule
		Selectors  []PageSelector
		Components []PageComponent
	}

	Declaration struct {
		Node
		Name      string
		Value     []token.Token
		Important bool
	}

	Selector interface {
		Node
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
		Selector
		Name string
	}

	ClassSelector struct {
		Selector
		Name string
	}

	AttributeSelector struct {
		Name      string
		Namespace *string
		Value     string
		Matcher   string
		Modifier  string
	}

	TypeSelector struct {
		Selector
		Name      string
		Namespace *string
	}

	PseudoSelector struct {
		Selector
		Name       string
		Functional bool
		Value      []token.Token
	}

	CompoundSelector struct {
		Selector
		Left  Selector
		Right Selector
	}

	RelativeSelector struct {
		Selector
		Combinator rune
		Left       Selector
		Right      Selector
	}

	MediaQuery struct {
		Node
		Type      string
		Qualifier string
		Condition MediaCondition
	}

	MediaCondition interface {
		Node
		VisitMediaCondition(MediaConditionVisitor)
	}

	MediaConditionVisitor struct {
		MediaOperation func(MediaOperation)
		MediaFeature   func(MediaFeature)
		MediaNegation  func(MediaNegation)
	}

	MediaOperation struct {
		MediaCondition
		Operator string
		Left     MediaCondition
		Right    MediaCondition
	}

	MediaFeature struct {
		MediaCondition
		Name  string
		Value MediaValue
	}

	MediaNegation struct {
		MediaCondition
		Condition MediaCondition
	}

	MediaValue interface {
		Node
		VisitMediaValue(MediaValueVisitor)
	}

	MediaValueVisitor struct {
		MediaValuePlain func(MediaValuePlain)
		MediaValueRange func(MediaValueRange)
	}

	MediaValuePlain struct {
		MediaValue
		Value token.Token
	}

	MediaValueRange struct {
		MediaValue
		LowerValue     token.Token
		LowerInclusive bool
		UpperValue     token.Token
		UpperInclusive bool
	}

	KeyframeBlock struct {
		Node
		Selector     float64
		Declarations []Declaration
	}

	SupportsCondition interface {
		Node
		VisitSupportsCondition(SupportsConditionVisitor)
	}

	SupportsConditionVisitor struct {
		SupportsOperation func(SupportsOperation)
		SupportsFeature   func(SupportsFeature)
		SupportsNegation  func(SupportsNegation)
	}

	SupportsOperation struct {
		SupportsCondition
		Operator string
		Left     SupportsCondition
		Right    SupportsCondition
	}

	SupportsFeature struct {
		SupportsCondition
		Declaration Declaration
	}

	SupportsNegation struct {
		SupportsCondition
		Condition SupportsCondition
	}

	PageSelector struct {
		Node
		Type    string
		Classes []string
	}

	PageComponent interface {
		Node
		VisitPageComponent(PageComponentVisitor)
	}

	PageComponentVisitor struct {
		PageDeclaration func(PageDeclaration)
		PageMargin      func(PageMargin)
	}

	PageDeclaration struct {
		PageComponent
		Declaration Declaration
	}

	PageMargin struct {
		PageComponent
		Name         string
		Declarations []Declaration
	}
)

func (r StyleRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r StyleRule) VisitRule(v RuleVisitor) { v.StyleRule(r) }

func (r ImportRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r ImportRule) VisitRule(v RuleVisitor) { v.ImportRule(r) }

func (r MediaRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r MediaRule) VisitRule(v RuleVisitor) { v.MediaRule(r) }

func (r KeyframesRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r KeyframesRule) VisitRule(v RuleVisitor) { v.KeyframesRule(r) }

func (r SupportsRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r SupportsRule) VisitRule(v RuleVisitor) { v.SupportsRule(r) }

func (r PageRule) VisitNode(v NodeVisitor) { v.Rule(r) }
func (r PageRule) VisitRule(v RuleVisitor) { v.PageRule(r) }

func (d Declaration) VisitNode(v NodeVisitor) { v.Declaration(d) }

func (s IdSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s IdSelector) VisitSelector(v SelectorVisitor) { v.IdSelector(s) }

func (s ClassSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s ClassSelector) VisitSelector(v SelectorVisitor) { v.ClassSelector(s) }

func (s AttributeSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s AttributeSelector) VisitSelector(v SelectorVisitor) { v.AttributeSelector(s) }

func (s TypeSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s TypeSelector) VisitSelector(v SelectorVisitor) { v.TypeSelector(s) }

func (s PseudoSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s PseudoSelector) VisitSelector(v SelectorVisitor) { v.PseudoSelector(s) }

func (s RelativeSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s RelativeSelector) VisitSelector(v SelectorVisitor) { v.RelativeSelector(s) }

func (s CompoundSelector) VisitNode(v NodeVisitor)         { v.Selector(s) }
func (s CompoundSelector) VisitSelector(v SelectorVisitor) { v.CompoundSelector(s) }

func (m MediaOperation) VisitMediaCondition(v MediaConditionVisitor) { v.MediaOperation(m) }

func (m MediaNegation) VisitMediaCondition(v MediaConditionVisitor) { v.MediaNegation(m) }

func (m MediaFeature) VisitMediaCondition(v MediaConditionVisitor) { v.MediaFeature(m) }

func (m MediaValuePlain) VisitMediaValue(v MediaValueVisitor) { v.MediaValuePlain(m) }

func (s SupportsOperation) VisitSupportsCondition(v SupportsConditionVisitor) { v.SupportsOperation(s) }

func (s SupportsNegation) VisitSupportsCondition(v SupportsConditionVisitor) { v.SupportsNegation(s) }

func (s SupportsFeature) VisitSupportsCondition(v SupportsConditionVisitor) { v.SupportsFeature(s) }

func (p PageDeclaration) VisitPageComponent(v PageComponentVisitor) { v.PageDeclaration(p) }
