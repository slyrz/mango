package markup

import (
	"fmt"
	"time"
)

type Writer interface {
	WriteTitle(name string)
	WriteDate(date time.Time)
	WritePart(name string)
	Write(format string, args ...interface{})

	Parts() map[string]string
	Order() []string

	Head() string
	Tail() string
}

func Save(writer Writer) {
	// Special sections we care about. We want to print them in this order,
	// no matter where they were defined.
	sections := []string{
		"name",
		"synopsis",
		"description",
		"options",
	}

	head := writer.Head()
	tail := writer.Tail()

	if len(head) > 0 {
		fmt.Print(head)
	}

	parts := writer.Parts()
	for _, name := range sections {
		if data, ok := parts[name]; ok {
			fmt.Print(data)
			// Delete from map so we don't print it again in the second pass.
			delete(parts, name)
		}
	}

	// Print the remaining parts in the original order.
	for _, name := range writer.Order() {
		if buffer, ok := parts[name]; ok {
			fmt.Print(buffer)
		}
	}

	if len(tail) > 0 {
		fmt.Print(tail)
	}
}
