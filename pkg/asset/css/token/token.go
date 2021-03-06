package token

type (
	Token interface {
		VisitToken(TokenVisitor)
	}

	TokenVisitor struct {
		Ident       func(Ident)
		Function    func(Function)
		AtKeyword   func(AtKeyword)
		Hash        func(Hash)
		String      func(String)
		Url         func(Url)
		Delim       func(Delim)
		Number      func(Number)
		Percentage  func(Percentage)
		Dimension   func(Dimension)
		Whitespace  func(Whitespace)
		Colon       func(Colon)
		Semicolon   func(Semicolon)
		Comma       func(Comma)
		OpenSquare  func(OpenSquare)
		CloseSquare func(CloseSquare)
		OpenParen   func(OpenParen)
		CloseParen  func(CloseParen)
		OpenCurly   func(OpenCurly)
		CloseCurly  func(CloseCurly)
	}

	Ident struct {
		Offset int
		Value  string
	}

	Function struct {
		Offset int
		Value  string
	}

	AtKeyword struct {
		Offset int
		Value  string
	}

	Hash struct {
		Offset int
		Value  string
		Id     bool
	}

	String struct {
		Offset int
		Mark   rune
		Value  string
	}

	Url struct {
		Offset int
		Value  string
	}

	Delim struct {
		Offset int
		Value  rune
	}

	Number struct {
		Offset  int
		Value   float64
		Integer bool
	}

	Percentage struct {
		Offset int
		Value  float64
	}

	Dimension struct {
		Offset  int
		Value   float64
		Integer bool
		Unit    string
	}

	Whitespace struct {
		Offset int
	}

	Colon struct {
		Offset int
	}

	Semicolon struct {
		Offset int
	}

	CloseCurly struct {
		Offset int
	}

	OpenCurly struct {
		Offset int
	}

	CloseParen struct {
		Offset int
	}

	OpenParen struct {
		Offset int
	}

	CloseSquare struct {
		Offset int
	}

	OpenSquare struct {
		Offset int
	}

	Comma struct {
		Offset int
	}
)

func (t Ident) VisitToken(v TokenVisitor) { v.Ident(t) }

func (t Function) VisitToken(v TokenVisitor) { v.Function(t) }

func (t AtKeyword) VisitToken(v TokenVisitor) { v.AtKeyword(t) }

func (t Hash) VisitToken(v TokenVisitor) { v.Hash(t) }

func (t String) VisitToken(v TokenVisitor) { v.String(t) }

func (t Url) VisitToken(v TokenVisitor) { v.Url(t) }

func (t Delim) VisitToken(v TokenVisitor) { v.Delim(t) }

func (t Number) VisitToken(v TokenVisitor) { v.Number(t) }

func (t Percentage) VisitToken(v TokenVisitor) { v.Percentage(t) }

func (t Dimension) VisitToken(v TokenVisitor) { v.Dimension(t) }

func (t Whitespace) VisitToken(v TokenVisitor) { v.Whitespace(t) }

func (t Colon) VisitToken(v TokenVisitor) { v.Colon(t) }

func (t Semicolon) VisitToken(v TokenVisitor) { v.Semicolon(t) }

func (t Comma) VisitToken(v TokenVisitor) { v.Comma(t) }

func (t OpenSquare) VisitToken(v TokenVisitor) { v.OpenSquare(t) }

func (t CloseSquare) VisitToken(v TokenVisitor) { v.CloseSquare(t) }

func (t OpenParen) VisitToken(v TokenVisitor) { v.OpenParen(t) }

func (t CloseParen) VisitToken(v TokenVisitor) { v.CloseParen(t) }

func (t OpenCurly) VisitToken(v TokenVisitor) { v.OpenCurly(t) }

func (t CloseCurly) VisitToken(v TokenVisitor) { v.CloseCurly(t) }

func Offset(token Token) int {
	switch t := token.(type) {
	case Ident:
		return t.Offset
	case Function:
		return t.Offset
	case AtKeyword:
		return t.Offset
	case Hash:
		return t.Offset
	case String:
		return t.Offset
	case Url:
		return t.Offset
	case Delim:
		return t.Offset
	case Number:
		return t.Offset
	case Percentage:
		return t.Offset
	case Dimension:
		return t.Offset
	case Whitespace:
		return t.Offset
	case Colon:
		return t.Offset
	case Semicolon:
		return t.Offset
	case Comma:
		return t.Offset
	case OpenSquare:
		return t.Offset
	case CloseSquare:
		return t.Offset
	case OpenParen:
		return t.Offset
	case CloseParen:
		return t.Offset
	case OpenCurly:
		return t.Offset
	case CloseCurly:
		return t.Offset
	}

	return -1
}
