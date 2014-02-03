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
	sections := []string{
		"name",
		"synopsis",
		"description",
		"options",
	}

	fmt.Println(writer.Head())

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

	fmt.Println(writer.Tail())

}
