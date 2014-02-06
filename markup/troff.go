package markup

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type TroffRenderer struct {
	Output Writer
}

func NewTroffRenderer(output Writer) *TroffRenderer {
	return &TroffRenderer{output}
}

func (tr *TroffRenderer) Group(n *Node) {
	// Don't indent the root node.
	if n.Parent != nil {
		tr.Output.Writeln(`.RS 4`)
	}
	RenderChilds(tr, n)
	if n.Parent != nil {
		tr.Output.Writeln(`.RE`)
	}
}

func (tr *TroffRenderer) Block(n *Node) {
	tr.Output.Writeln(`.RS 4`)
	tr.Output.Writeln(`.nf`)
	RenderChilds(tr, n)
	tr.Output.Writeln(`.fi`)
	tr.Output.Writeln(`.RE`)
}

func (tr *TroffRenderer) List(n *Node) {
	RenderChilds(tr, n)
}

func (tr *TroffRenderer) ListItem(n *Node) {
	tr.Output.Writeln(`.TP`)
	if n.Text == "*" {
		tr.Output.Writeln(`\(bu`)
	} else {
		tr.Output.Writeln(`.B "%s"`, n.Text)
	}
	RenderChilds(tr, n)
}

func (tr *TroffRenderer) Section(text string) {
	tr.Output.WritePart(text)
	tr.Output.Writeln(`.SH "%s"`, strings.ToUpper(text))
}

func (tr *TroffRenderer) Space() {
	return
}

func (tr *TroffRenderer) Break() {
	tr.Output.Writeln(`.PP`)
}

func (tr *TroffRenderer) Text(text string) {
	tr.Output.Writeln(strings.TrimSpace(text))
}

func (tr *TroffRenderer) TextBold(text string) {
	tr.Output.Writeln(`.B "%s"`, text)
}

func (tr *TroffRenderer) TextUnderline(text string) {
	tr.Output.Writeln(`.I "%s"`, text)
}

type TroffWriter struct {
	title  string
	date   time.Time
	parts  map[string]*bytes.Buffer
	active *bytes.Buffer
	order  []string
}

func (w *TroffWriter) WriteTitle(name string) {
	w.title = name
}

func (w *TroffWriter) WriteDate(date time.Time) {
	w.date = date
}

func (w *TroffWriter) WritePart(name string) {
	key := strings.ToLower(name)
	w.order = append(w.order, key)
	w.active = new(bytes.Buffer)
	w.parts[key] = w.active
}

func (w *TroffWriter) writeFmt(newline bool, format string, args ...interface{}) {
	if w.active == nil {
		panic("write called but no part active")
	}
	fmt.Fprintf(w.active, format, args...)
	if newline {
		fmt.Fprintln(w.active)
	}
}

func (w *TroffWriter) Write(format string, args ...interface{}) {
	w.writeFmt(false, format, args...)
}

func (w *TroffWriter) Writeln(format string, args ...interface{}) {
	w.writeFmt(true, format, args...)
}

func (w *TroffWriter) Parts() map[string]string {
	result := make(map[string]string)
	for key, buffer := range w.parts {
		result[key] = buffer.String()
	}
	return result
}

func (w *TroffWriter) Order() []string {
	return w.order
}

func (w *TroffWriter) Head() string {
	// Generates a Linux style man page title line:
	//   - Title of man page in all caps
	//   - Section number
	//   - Date in YYYY-MM-DD format (position footer, middle)
	//   - Source of the command (position footer, left)
	//   - Title of the manual (position header, center)
	titleUp := strings.ToUpper(w.title)
	titleHi := strings.Title(w.title)
	dateStr := fmt.Sprintf("%d-%02d-%02d", w.date.Year(), w.date.Month(), w.date.Day())

	return fmt.Sprintf(`.TH "%s" 1 "%s" "%s" "%s Manual"`+"\n", titleUp, dateStr, titleHi, titleHi)
}

func (w *TroffWriter) Tail() string {
	return ""
}

func NewTroffWriter() *TroffWriter {
	result := new(TroffWriter)
	result.parts = make(map[string]*bytes.Buffer)
	result.order = make([]string, 0)
	result.title = ""
	result.date = time.Now()
	return result
}
