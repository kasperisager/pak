package runes

import (
	"testing"
)

func TestLines(t *testing.T) {
	lines := Lines([]rune("Foo"))

	for _, line := range lines {
		t.Logf("%d %#v", line.Offset, string(line.Value))
	}

	t.Logf("%#v", LineAtOffset(lines, 0))
}
