package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var options = struct {
	Path   string
	Name   string
	Plain  bool
	Stdout bool
}{}

func init() {
	flag.StringVar(&options.Path, "dir", "", "output directory")
	flag.StringVar(&options.Name, "name", "", "command name")
	flag.BoolVar(&options.Plain, "plain", false, "plain text comments")
	flag.BoolVar(&options.Stdout, "stdout", false, "write output to stdout")
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

		if options.Stdout {
			fmt.Println(text)
		} else {
			dst, err := os.Create(path.Join(options.Path, file.Name))
			if err != nil {
				panic(err)
			}
			dst.Write([]byte(text))
			dst.Close()
		}
	}
}
