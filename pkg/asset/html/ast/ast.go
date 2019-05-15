package ast

type (
	Node interface {
		VisitNode(NodeVisitor)
	}

	NodeVisitor struct {
		Element   func(Element)
		Attribute func(Attribute)
		Text      func(Text)
	}

	Element struct {
		Name       string
		Attributes []Attribute
		Children   []Node
	}

	Attribute struct {
		Name  string
		Value string
	}

	Text struct {
		Data string
	}
)

func (e Element) VisitNode(v NodeVisitor) {
	v.Element(e)
}

func (e Element) Attribute(name string) (string, bool) {
	for _, attribute := range e.Attributes {
		if name == attribute.Name {
			return attribute.Value, true
		}
	}

	return "", false
}

func (e Element) IsVoid() bool {
	switch e.Name {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr":
		return true
	}

	return false
}

func (a Attribute) VisitNode(v NodeVisitor) {
	v.Attribute(a)
}

func (t Text) VisitNode(v NodeVisitor) {
	v.Text(t)
}
