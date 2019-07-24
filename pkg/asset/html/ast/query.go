package ast

type Query func(*Element) bool

func (q Query) And(p Query) Query {
	return func(element *Element) bool {
		return q(element) && p(element)
	}
}

func (q Query) Or(p Query) Query {
	return func(element *Element) bool {
		return q(element) || p(element)
	}
}

func ByName(name string) Query {
	return func(element *Element) bool {
		return element.Name == name
	}
}

func ByAttribute(name string, value string) Query {
	return func(element *Element) bool {
		found, ok := element.Attribute(name)

		if ok {
			if found == value {
				return true
			}
		}

		return false
	}
}
