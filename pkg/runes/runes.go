package runes

const EOF = -1

type Runes []rune

func (runes Runes) Peek(n int) rune {
	if len(runes) < n {
		return EOF
	}

	return runes[n-1]
}
