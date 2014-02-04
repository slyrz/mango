package markup

import (
	"bufio"
	"io"
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

func Save(writer Writer, out io.Writer) error {
	// Special sections we care about. We want to print them in this order,
	// no matter where they were defined.
	sections := []string{
		"name",
		"synopsis",
		"description",
		"options",
	}

	bufOut := bufio.NewWriter(out)
	defer bufOut.Flush()

	head := writer.Head()
	tail := writer.Tail()

	if _, err := bufOut.WriteString(head); err != nil {
		return err
	}
	parts := writer.Parts()
	for _, name := range sections {
		if data, ok := parts[name]; ok {
			if _, err := bufOut.WriteString(data); err != nil {
				return err
			}
			// Delete from map so we don't print it again in the second pass.
			delete(parts, name)
		}
	}

	// Print the remaining parts in the original order.
	for _, name := range writer.Order() {
		if data, ok := parts[name]; ok {
			if _, err := bufOut.WriteString(data); err != nil {
				return err
			}
		}
	}

	if _, err := bufOut.WriteString(tail); err != nil {
		return err
	}

	return nil
}
