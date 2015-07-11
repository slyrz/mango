package main

import (
	"flag"
	_ "fmt"
	"os"
	"os/exec"
)

var options = struct {
	Output  string
	Name    string
	Plain   bool
	Preview bool
}{}

func init() {
	flag.StringVar(&options.Output, "output", "", "write output to file")
	flag.StringVar(&options.Name, "name", "", "set command name")
	flag.BoolVar(&options.Plain, "plain", false, "treat comments as plain text")
	flag.BoolVar(&options.Preview, "preview", false, "preview output in man")
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
		if options.Name != "" {
			file.Name = options.Name
		}
		text, err := builder.Build(file)
		if err != nil {
			panic(err)
		}
		if options.Preview {
			cmd := exec.Command("groff", "-Wall", "-mtty-char", "-mandoc", "-Tascii")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			inp, err := cmd.StdinPipe()
			if err != nil {
				panic(err)
			}
			if err := cmd.Start(); err != nil {
				panic(err)
			}
			inp.Write([]byte(text))
			inp.Close()
			cmd.Wait()
		} else {
			if options.Output == "" {
				os.Stdout.Write([]byte(text))
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
}
