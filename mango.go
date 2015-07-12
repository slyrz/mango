// mango - generate manual pages from Go source code
//
// Description:
//
// mango generates manual pages from the source code of Go commands.
// It aims to generate full-fledged manual pages soley based on the comments
// and flag function calls found inside Go source code.
//
// See Also:
//
// man(1), man-pages(7)
package main

import (
	"flag"
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
	// Write the manual page to file.
	flag.StringVar(&options.Output, "output", "", "write to `file`")
	// Set the manual page title to name.
	flag.StringVar(&options.Name, "title", "", "set title to `name`")
	// Ignore markup inside comments.
	flag.BoolVar(&options.Plain, "plain", false, "treat comments as plain text")
	// Preview the manual page with the man command.
	flag.BoolVar(&options.Preview, "preview", false, "preview with man")
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
