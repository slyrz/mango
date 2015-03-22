package main

import (
	"bytes"
	"fmt"
	"strings"
)

type Writer interface {
	Block(string)
	Break(string)
	List(string)
	Section(string)
	Text(string)
	TextBold(string)
	TextUnderline(string)
	Option(string, string, string)

	Done() string
}

var manSections = []string{
	"name",
	"synopsis",
	"description",
	"options",
	"exit status",
}

type Troff struct {
	sections map[string]string
	order    []string
	active   string
	buffer   bytes.Buffer
}

func NewTroffWriter() *Troff {
	return &Troff{
		sections: make(map[string]string),
	}
}

func (tr *Troff) Done() string {
	if tr.buffer.Len() > 0 {
		tr.sections[tr.active] = tr.buffer.String()
	}
	tr.buffer.Reset()

	tr.writeln(`.TH "%s" 1 "%s" "%s" "%s Manual"`, "TEST", "1999-12-31", "TEST", "TEST")

	// At first, render special Manpage sections in their usual order.
	for _, section := range manSections {
		if output, ok := tr.sections[section]; ok {
			tr.write(output)
			delete(tr.sections, section)
		}
	}

	// Now render the remaining sections in the order they appeard in the
	// source file.
	for _, section := range tr.order {
		if output, ok := tr.sections[section]; ok {
			tr.write(output)
		}
	}

	return tr.buffer.String()
}

func (tr *Troff) write(format string, args ...interface{}) {
	fmt.Fprintf(&tr.buffer, format, args...)
}

func (tr *Troff) writeln(format string, args ...interface{}) {
	fmt.Fprintf(&tr.buffer, format+"\n", args...)
}

func (tr *Troff) Break(text string) {
	tr.writeln(".PP")
}

func (tr *Troff) Block(text string) {
	tr.writeln(".RS 4")
	tr.writeln(".nf")
	tr.writeln(text)
	tr.writeln(".fi")
	tr.writeln(".RE")
}

func (tr *Troff) List(text string) {
	tr.writeln(".TP")
	if text == "*" {
		tr.writeln(`\(bu`)
	} else {
		tr.writeln(`.B "%s"`, text)
	}
}

func (tr *Troff) Section(text string) {
	if tr.buffer.Len() > 0 {
		tr.sections[tr.active] = tr.buffer.String()
	}
	tr.buffer.Reset()
	tr.active = strings.ToLower(text)
	tr.order = append(tr.order, tr.active)

	tr.writeln(`.SH "%s"`, strings.ToUpper(text))
}

func (tr *Troff) Text(text string) {
	text = strings.TrimSpace(text)
	if text != "" {
		tr.writeln(text)
	}
}

func (tr *Troff) TextBold(text string) {
	tr.writeln(`.B "%s"`, strings.TrimSpace(text))
}

func (tr *Troff) TextUnderline(text string) {
	tr.writeln(`.I "%s"`, strings.TrimSpace(text))
}

func (tr *Troff) Option(name, short, param string) {
	tr.writeln(".TP")
	if short != "" {
		tr.write(`.B \-%s -%s`, short, name)
	} else {
		tr.write(`.B \-%s`, name)
	}
	if param != "" {
		tr.write(` \fI%s\fR`, param)
	}
	tr.writeln("")

}
