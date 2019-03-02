package runes

import (
	"sort"
)

type Line struct {
	Offset int
	Value  []rune
}

func Lines(runes []rune) (lines []Line) {
	start := 0

	for start < len(runes) {
		end := start

		for i := start; i < len(runes); i++ {
			end++

			if runes[i] == '\n' {
				break
			}

			if runes[i] == '\r' {
				if i+1 < len(runes) && runes[i+1] == '\n' {
					end++
				}

				break
			}
		}

		lines = append(lines, Line{
			Offset: start,
			Value:  runes[start:end],
		})

		start = end
	}

	return lines
}

func LineAtOffset(lines []Line, offset int) Line {
	i := sort.Search(len(lines), func(i int) bool {
		return lines[i].Offset > offset
	})

	return lines[i-1]
}
