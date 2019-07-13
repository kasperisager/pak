package token

type (
	Token interface {
		VisitToken(TokenVisitor)
	}

	TokenVisitor struct {
		Whitespace     func(Whitespace)
		LineTerminator func(LineTerminator)
		Keyword        func(Keyword)
		Identifier     func(Identifier)
		Punctuator     func(Punctuator)
		NumericLiteral func(NumericLiteral)
		StringLiteral  func(StringLiteral)
		BooleanLiteral func(BooleanLiteral)
		NullLiteral    func(NullLiteral)
		Template       func(Template)
		TemplateHead   func(TemplateHead)
		TemplateMiddle func(TemplateMiddle)
		TemplateTail   func(TemplateTail)
	}

	Whitespace struct {
		Offset int
	}

	LineTerminator struct {
		Offset int
	}

	Keyword struct {
		Offset int
		Value  string
	}

	Identifier struct {
		Offset int
		Value  string
	}

	Punctuator struct {
		Offset int
		Value  string
	}

	NumericLiteral struct {
		Offset int
		Value  float64
	}

	StringLiteral struct {
		Offset int
		Value  string
	}

	NullLiteral struct {
		Offset int
	}

	BooleanLiteral struct {
		Offset int
		Value  bool
	}

	Template struct {
		Offset int
		Value  string
	}

	TemplateHead struct {
		Offset int
		Value  string
	}

	TemplateMiddle struct {
		Offset int
		Value  string
	}

	TemplateTail struct {
		Offset int
		Value  string
	}
)

func (t Whitespace) VisitToken(v TokenVisitor) { v.Whitespace(t) }

func (t Keyword) VisitToken(v TokenVisitor) { v.Keyword(t) }

func (t Identifier) VisitToken(v TokenVisitor) { v.Identifier(t) }

func (t Punctuator) VisitToken(v TokenVisitor) { v.Punctuator(t) }

func (t StringLiteral) VisitToken(v TokenVisitor) { v.StringLiteral(t) }

func (t BooleanLiteral) VisitToken(v TokenVisitor) { v.BooleanLiteral(t) }

func (t NullLiteral) VisitToken(v TokenVisitor) { v.NullLiteral(t) }
