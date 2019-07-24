package ast

type Iterator func() (*Element, bool)

func (i Iterator) Next() (*Element, bool) {
	return i()
}
