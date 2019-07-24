package ast

type (
	Node interface {
		VisitNode(NodeVisitor)
	}

	NodeVisitor struct {
		Element func(*Element)
		Text    func(*Text)
	}

	Element struct {
		Name       string
		Attributes []*Attribute
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

func (e *Element) VisitNode(v NodeVisitor) { v.Element(e) }
func (t *Text) VisitNode(v NodeVisitor)    { v.Text(t) }

func (e *Element) Attribute(name string) (string, bool) {
	for _, attribute := range e.Attributes {
		if name == attribute.Name {
			return attribute.Value, true
		}
	}

	return "", false
}

func (e *Element) Text() string {
	var text string

	for _, child := range e.Children {
		switch child := child.(type) {
		case *Element:
			text += child.Text()

		case *Text:
			text += child.Data
		}
	}

	return text
}

func (e *Element) IsVoid() bool {
	switch e.Name {
	case
		"area",
		"base",
		"br",
		"col",
		"embed",
		"hr",
		"img",
		"input",
		"link",
		"meta",
		"param",
		"source",
		"track",
		"wbr":
		return true
	}

	return false
}

func (e *Element) Walk() Iterator {
	queue := []*Element{e}

	return func() (*Element, bool) {
		if len(queue) > 0 {
			var element *Element

			element, queue = queue[0], queue[1:]

			for _, child := range element.Children {
				switch child := child.(type) {
				case *Element:
					queue = append(queue, child)
				}
			}

			return element, true
		}

		return nil, false
	}
}

func (e *Element) Find(query Query) Iterator {
	it := e.Walk()

	return func() (*Element, bool) {
		for {
			element, ok := it.Next()

			if !ok {
				break
			}

			if query(element) {
				return element, true
			}
		}

		return nil, false
	}
}
