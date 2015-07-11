package main

import (
	"flag"
	"fmt"
	"os"
)

var options = struct {
	Output string
	Name   string
	Plain  bool
}{}

func init() {
	flag.StringVar(&options.Output, "output", "", "write output to file")
	flag.StringVar(&options.Name, "name", "", "command name")
	flag.BoolVar(&options.Plain, "plain", false, "plain text comments")
}

func getReader() Reader {
	if options.Plain {
		return NewPlainReader()
	} else {
		return NewMarkupReader()
	}
}

func getWriter() Writer {
	return NewTroffWriter()
}

func main() {
	flag.Parse()

	builder := NewBuilder(getReader(), getWriter())
	for _, arg := range flag.Args() {
		file, err := NewFile(arg)
		if err != nil {
			panic(err)
		}

		text, err := builder.Build(file)
		if err != nil {
			panic(err)
		}
		if options.Output == "" {
			fmt.Println(text)
		} else {
			dst, err := os.Create(options.Output)
			if err != nil {
				panic(err)
			}
			dst.Write([]byte(text))
			dst.Close()
		}
	}
}
