// mango - generate man pages from the source of your Go commands
//
// Description:
//
// TODO...
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
	b.Writer.WriteDate(b.File.Time)
	b.feedDocumentation()
	b.feedSynopsis()
	b.feedOptions()
	return nil
}

func (b *Builder) feedDocumentation() {
	tokens, err := b.Tokenizer.TokenizeString(b.File.Doc)
	if err != nil {
		return
	}

	b.Renderer.Section("Name")
	markup.Render(b.Renderer, b.Parser.Parse(tokens))
}

func (b *Builder) feedSynopsis() {
	b.Renderer.Section("Synopsis")
	b.Renderer.Text(b.File.Name)
	if len(b.File.Options) > 0 {
		b.Renderer.TextUnderline("[option]")
		b.Renderer.Text("... ")
	}
	b.Renderer.TextUnderline("[args]")
	b.Renderer.Text("... ")
	b.Renderer.Break()
}

func (b *Builder) feedOptions() {
	if len(b.File.Options) == 0 {
		return
	}

	b.Renderer.Section("Options")
	for _, opt := range b.File.Options {
		textHead := ""
		textBody := ""

		if len(opt.Short) > 0 {
			textHead = fmt.Sprintf("-%s, -%s", opt.Short, opt.Name)
		} else {
			textHead = fmt.Sprintf("-%s", opt.Name)
		}

		if len(opt.Doc) > 0 {
			textBody = opt.Doc
		} else {
			textBody = opt.Usage
		}

		// Tokenize body text. We haven't written anything yet, so if Tokenize
		// function fails, the document stays unchanged and we try to parse the
		// next option.
		tokens, err := b.Tokenizer.TokenizeString(textBody)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning:")
			continue
		}

		b.Renderer.TextBold(textHead)
		if len(tokens) > 0 {
			markup.Render(b.Renderer, b.Parser.ParsePart(tokens))
		}
		b.Renderer.Break()
	}
}

func (b *Builder) Save(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return markup.Save(b.Writer, file)
}

func main() {
	for _, srcPath := range os.Args[1:] {
		builder := NewBuilder()

		if err := builder.Load(srcPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not open file '%s': %s\n", srcPath, err)
			continue
		}

		dstPath := fmt.Sprintf("%s.1", builder.File.Name)
		if err := builder.Save(dstPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not save '%s': %s\n", dstPath, err)
			continue
		}

		fmt.Printf("%s -> %s\n", srcPath, dstPath)
	}
}
