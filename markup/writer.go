package markup

import (
	"bufio"
	"io"
	"time"
)

var (
	// The array below contains a few conventional or suggested sections, that
	// should be placed at the start of the manual page in the order shown in
	// the list. We check if the documentation contains some of these sections
	// and print the ones we found first.
	sections = []string{
		"name",
		"synopsis",
		"description",
		"options",
		"exit status",
	}
)

type Writer interface {
	WriteTitle(name string)
	WriteDate(date time.Time)
	WritePart(name string)

	Write(format string, args ...interface{})
	Writeln(format string, args ...interface{})

	Parts() map[string]string
	Order() []string

	Head() string
	Tail() string
}

func Save(writer Writer, out io.Writer) error {
	bufOut := bufio.NewWriter(out)
	defer bufOut.Flush()

	head := writer.Head()
	tail := writer.Tail()

	// Head contains stuff the writer needs to put at the document start.
	if _, err := bufOut.WriteString(head); err != nil {
		return err
	}
	// Print the conventional manual page sections.
	parts := writer.Parts()
	for _, name := range sections {
		if data, ok := parts[name]; ok {
			if _, err := bufOut.WriteString(data); err != nil {
				return err
			}
			// Delete section from map so we don't print it again
			// in the second pass.
			delete(parts, name)
		}
	}
	// After we printed the conventional sections, we print the user defined
	// sections now and we keep their order.
	for _, name := range writer.Order() {
		if data, ok := parts[name]; ok {
			if _, err := bufOut.WriteString(data); err != nil {
				return err
			}
		}
	}
	// Tail contains stuff the writer needs to put at the document end.
	if _, err := bufOut.WriteString(tail); err != nil {
		return err
	}
	return nil
}
