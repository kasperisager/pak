package token

type (
	Token interface {
		VisitToken(TokenVisitor)
	}

	TokenVisitor struct {
		Newline        func(Newline)
		Keyword        func(Keyword)
		Identifier     func(Identifier)
		Punctuator     func(Punctuator)
		Number         func(Number)
		String         func(String)
		Boolean        func(Boolean)
		Null           func(Null)
		Template       func(Template)
		TemplateHead   func(TemplateHead)
		TemplateMiddle func(TemplateMiddle)
		TemplateTail   func(TemplateTail)
	}

	Newline struct {
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

	Number struct {
		Offset int
		Value  float64
	}

	String struct {
		Offset int
		Value  string
	}

	Null struct {
		Offset int
	}

	Boolean struct {
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

func (t Keyword) VisitToken(v TokenVisitor)    { v.Keyword(t) }
func (t Identifier) VisitToken(v TokenVisitor) { v.Identifier(t) }
func (t Punctuator) VisitToken(v TokenVisitor) { v.Punctuator(t) }
func (t Number) VisitToken(v TokenVisitor)     { v.Number(t) }
func (t String) VisitToken(v TokenVisitor)     { v.String(t) }
func (t Boolean) VisitToken(v TokenVisitor)    { v.Boolean(t) }
func (t Null) VisitToken(v TokenVisitor)       { v.Null(t) }
