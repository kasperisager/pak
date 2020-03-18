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
		StyleRule     func(*StyleRule)
		ImportRule    func(*ImportRule)
		MediaRule     func(*MediaRule)
		FontFaceRule  func(*FontFaceRule)
		KeyframesRule func(*KeyframesRule)
		SupportsRule  func(*SupportsRule)
		PageRule      func(*PageRule)
	}

	StyleRule struct {
		Selectors    []Selector
		Declarations []*Declaration
	}

	ImportRule struct {
		URL        *url.URL
		Conditions []*MediaQuery
	}

	MediaRule struct {
		Conditions []*MediaQuery
		StyleSheet *StyleSheet
	}

	FontFaceRule struct {
		Declarations []*Declaration
	}

	KeyframesRule struct {
		Prefix string
		Name   string
		Blocks []*KeyframeBlock
	}

	SupportsRule struct {
		Condition  SupportsCondition
		StyleSheet *StyleSheet
	}

	PageRule struct {
		Selectors  []*PageSelector
		Components []PageComponent
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
		IdSelector        func(*IdSelector)
		ClassSelector     func(*ClassSelector)
		AttributeSelector func(*AttributeSelector)
		TypeSelector      func(*TypeSelector)
		PseudoSelector    func(*PseudoSelector)
		CompoundSelector  func(*CompoundSelector)
		ComplexSelector   func(*ComplexSelector)
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
		Matcher   string
		Modifier  string
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

	ComplexSelector struct {
		Combinator rune
		Left       Selector
		Right      Selector
	}

	MediaQuery struct {
		Type      string
		Qualifier string
		Condition MediaCondition
	}

	MediaCondition interface {
		VisitMediaCondition(MediaConditionVisitor)
	}

	MediaConditionVisitor struct {
		MediaOperation func(*MediaOperation)
		MediaFeature   func(*MediaFeature)
		MediaNegation  func(*MediaNegation)
	}

	MediaOperation struct {
		Operator string
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
		MediaValuePlain func(*MediaValuePlain)
		MediaValueRange func(*MediaValueRange)
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

	KeyframeBlock struct {
		Selector     float64
		Declarations []*Declaration
	}

	SupportsCondition interface {
		VisitSupportsCondition(SupportsConditionVisitor)
	}

	SupportsConditionVisitor struct {
		SupportsOperation func(*SupportsOperation)
		SupportsFeature   func(*SupportsFeature)
		SupportsNegation  func(*SupportsNegation)
	}

	SupportsOperation struct {
		Operator string
		Left     SupportsCondition
		Right    SupportsCondition
	}

	SupportsFeature struct {
		Declaration *Declaration
	}

	SupportsNegation struct {
		Condition SupportsCondition
	}

	PageSelector struct {
		Type    string
		Classes []string
	}

	PageComponent interface {
		VisitPageComponent(PageComponentVisitor)
	}

	PageComponentVisitor struct {
		PageDeclaration func(*PageDeclaration)
		PageMargin      func(*PageMargin)
	}

	PageDeclaration struct {
		Declaration *Declaration
	}

	PageMargin struct {
		Name         string
		Declarations []*Declaration
	}
)

func (r *StyleRule) VisitRule(v RuleVisitor)     { v.StyleRule(r) }
func (r *ImportRule) VisitRule(v RuleVisitor)    { v.ImportRule(r) }
func (r *MediaRule) VisitRule(v RuleVisitor)     { v.MediaRule(r) }
func (r *FontFaceRule) VisitRule(v RuleVisitor)  { v.FontFaceRule(r) }
func (r *KeyframesRule) VisitRule(v RuleVisitor) { v.KeyframesRule(r) }
func (r *SupportsRule) VisitRule(v RuleVisitor)  { v.SupportsRule(r) }
func (r *PageRule) VisitRule(v RuleVisitor)      { v.PageRule(r) }

func (s *IdSelector) VisitSelector(v SelectorVisitor)        { v.IdSelector(s) }
func (s *ClassSelector) VisitSelector(v SelectorVisitor)     { v.ClassSelector(s) }
func (s *AttributeSelector) VisitSelector(v SelectorVisitor) { v.AttributeSelector(s) }
func (s *TypeSelector) VisitSelector(v SelectorVisitor)      { v.TypeSelector(s) }
func (s *PseudoSelector) VisitSelector(v SelectorVisitor)    { v.PseudoSelector(s) }
func (s *ComplexSelector) VisitSelector(v SelectorVisitor)   { v.ComplexSelector(s) }
func (s *CompoundSelector) VisitSelector(v SelectorVisitor)  { v.CompoundSelector(s) }

func (m *MediaOperation) VisitMediaCondition(v MediaConditionVisitor) { v.MediaOperation(m) }
func (m *MediaNegation) VisitMediaCondition(v MediaConditionVisitor)  { v.MediaNegation(m) }
func (m *MediaFeature) VisitMediaCondition(v MediaConditionVisitor)   { v.MediaFeature(m) }

func (m *MediaValuePlain) VisitMediaValue(v MediaValueVisitor) { v.MediaValuePlain(m) }

func (s *SupportsOperation) VisitSupportsCondition(v SupportsConditionVisitor) { v.SupportsOperation(s) }
func (s *SupportsNegation) VisitSupportsCondition(v SupportsConditionVisitor)  { v.SupportsNegation(s) }
func (s *SupportsFeature) VisitSupportsCondition(v SupportsConditionVisitor)   { v.SupportsFeature(s) }

func (p *PageDeclaration) VisitPageComponent(v PageComponentVisitor) { v.PageDeclaration(p) }
