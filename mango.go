// mango - generate man pages from the source of your Go commands
//
// Descrition
// ----------
// Lorem ipsum dolor sit amet, *consetetur* sadipscing elitr, sed diam nonumy
// eirmod tempor invidunt ut labore et _dolore_ magna aliquyam erat, sed diam
// voluptua.
//
//  	echo "Hello World"
//
// At vero eos et accusam et justo duo dolores et ea rebum. Stet clita
// kasd gubergren, no sea takimata sanctus est.
//
// Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy
// eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam
// voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet
// clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
//
// Return Codes
// ------------
// 0) Success
// 1) EOF
// *) Let's try and see end
//    what happens here.
//
package main

import (
	"./markup"
	"./source"
	"fmt"
	"os"
)

type Builder struct {
	File      *source.File
	Tokenizer *markup.Tokenizer
	Parser    *markup.Parser
	Root      *markup.Node
	Renderer  markup.Renderer
	Writer    markup.Writer
}

func NewBuilder() *Builder {
	result := new(Builder)
	result.Tokenizer = markup.NewTokenizer()
	result.Parser = markup.NewParser()
	result.Writer = markup.NewTroffWriter()
	result.Renderer = markup.NewTroffRenderer(result.Writer)

	return result
}

func (b *Builder) Load(path string) error {
	file, err := source.NewFile(path)
	if err != nil {
		return err
	}
	b.File = file

	b.Writer.WriteTitle(b.File.Name)

	tokens, err := b.Tokenizer.TokenizeString(b.File.Doc)
	if err != nil {
		return err
	}

	b.Root = b.Parser.Parse(tokens)
	b.Renderer.Section("Name")
	markup.Render(b.Renderer, b.Root)

	b.feedSynopsis()
	b.feedOptions()

	return nil
}

func (b *Builder) feedSynopsis() {
	b.Renderer.Section("Synopsis")
	b.Renderer.Text(b.File.Name)
	b.Renderer.TextUnderline("[option]")
	b.Renderer.Text("... ")
	b.Renderer.TextUnderline("[args]")
	b.Renderer.Text("... ")
	b.Renderer.Break()
}

func (b *Builder) feedOptions() {
	b.Renderer.Section("options")

	for _, opt := range b.File.Options {
		if len(opt.Short) > 0 {
			b.Renderer.TextBold(fmt.Sprintf("-%s, ", opt.Short))
		}
		b.Renderer.TextBold(fmt.Sprintf("-%s", opt.Name))

		text := ""
		if len(opt.Doc) > 0 {
			text = opt.Doc
		} else {
			text = opt.Usage
		}

		tokens, err := b.Tokenizer.TokenizeString(text)
		if err != nil {
			panic(err)
		}

		node := b.Parser.Parse(tokens)
		node.Parent = b.Root
		markup.Render(b.Renderer, node)
		b.Renderer.Break()
	}
}

func (b *Builder) Save(path string) {
	markup.Save(b.Writer)
}

func main() {
	for _, filename := range os.Args[1:] {
		builder := NewBuilder()
		builder.Load(filename)
		builder.Save("")
	}
}
