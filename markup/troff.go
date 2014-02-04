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
		tr.Output.Write(".RS 4\n")
	}
	RenderChilds(tr, n)
	if n.Parent != nil {
		tr.Output.Write(".RE\n")
	}
}

func (tr *TroffRenderer) Block(n *Node) {
	tr.Output.Write(".RS 4\n")
	tr.Output.Write(".nf\n")
	RenderChilds(tr, n)
	tr.Output.Write(".fi\n")
	tr.Output.Write(".RE\n")
}

func (tr *TroffRenderer) List(n *Node) {
	RenderChilds(tr, n)
}

func (tr *TroffRenderer) ListItem(n *Node) {
	tr.Output.Write(".TP\n")
	if n.Text == "*" {
		tr.Output.Write("\\(bu\n")
	} else {
		tr.Output.Write(".B \"%s\"\n", n.Text)
	}
	RenderChilds(tr, n)
}

func (tr *TroffRenderer) Section(text string) {
	tr.Output.WritePart(text)
	tr.Output.Write(".SH \"%s\"\n", strings.ToUpper(text))
}

func (tr *TroffRenderer) Space() {
	return
}

func (tr *TroffRenderer) Break() {
	tr.Output.Write(".PP\n")
}

func (tr *TroffRenderer) Text(text string) {
	tr.Output.Write("%s\n", strings.TrimSpace(text))
}

func (tr *TroffRenderer) TextBold(text string) {
	tr.Output.Write(".B \"%s\"\n", text)
}

func (tr *TroffRenderer) TextUnderline(text string) {
	tr.Output.Write(".I \"%s\"\n", text)
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

func (w *TroffWriter) Write(format string, args ...interface{}) {
	if w.active == nil {
		panic("write called but no part active")
	}
	fmt.Fprintf(w.active, format, args...)
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
	titleLo := strings.ToLower(w.title)
	titleHi := strings.Title(w.title)
	dateStr := fmt.Sprintf("%02d/%02d/%d", w.date.Day(), w.date.Month(), w.date.Year())

	return fmt.Sprintf(".TH \"%s\" \"1\" \"\" \"%s\" \"%s Manual\"\n", titleLo, dateStr, titleHi)
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
