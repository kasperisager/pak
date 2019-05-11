package token

type (
	Token interface {
		VisitToken(TokenVisitor)
	}

	TokenVisitor struct {
		DocumentType func(DocumentType)
		StartTag func(StartTag)
		EndTag func(EndTag)
		Character func(Character)
	}

	DocumentType struct {
		Offset int
	}

	Attribute struct {
		Offset int
		Name  string
		Value string
	}

	StartTag struct {
		Offset int
		Name string
		Attributes []Attribute
		Closed bool
	}

	EndTag struct {
		Offset int
		Name string
	}

	Character struct {
		Offset int
		Data rune
	}
)

func (t DocumentType) VisitToken(v TokenVisitor) { v.DocumentType(t) }

func (t StartTag) VisitToken(v TokenVisitor) { v.StartTag(t) }

func (t EndTag) VisitToken(v TokenVisitor) { v.EndTag(t) }

func (t Character) VisitToken(v TokenVisitor) { v.Character(t) }
