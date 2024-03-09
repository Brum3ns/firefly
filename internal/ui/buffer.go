package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type buffer struct {
	data   [10]string
	head   int
	tail   int
	length int
}

func (b *buffer) Append(item string) {
	if b.length >= 10 {
		b.data[b.head] = item
		b.head = (b.head + 1) % 10
		b.tail = (b.tail + 1) % 10

	} else {
		b.data[b.head] = item
		b.head = (b.head + 1) % 10
		b.length++
	}
}

func (b *buffer) Print(style ...lipgloss.Style) {
	for i := 0; i < b.length; i++ {
		index := (b.tail + 1) % 10

		if len(style) > 0 {
			//style[0]
		} else {
			fmt.Println("%s", b.data[index])
		}
	}
}
