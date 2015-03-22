package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"time"
)

var ErrFileType = errors.New("not a Go file")

// File represents a parsed '.go' source file.
type File struct {
	Path    string    // Path to file.
	Name    string    // Name of command.
	Time    time.Time // Modification time.
	Options []*Option // Options found in file.
	Doc     string    // Comment preceding the "package" keyword.
}

func splitExt(s string) (string, string) {
	i := len(s) - len(path.Ext(s))
	return s[:i], s[i:]
}

func NewFile(path string) (*File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	name, ext := splitExt(info.Name())
	if ext != ".go" {
		return nil, ErrFileType
	}
	file := &File{
		Path: path,
		Name: name,
		Time: info.ModTime(),
	}
	return file, file.parseOptions()
}

// parseOptions transforms all flag package calls to Options.
func (f *File) parseOptions() error {
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, f.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// The last comment group before a package declaration must contain the
	// command description.
	packageLine := 2
	if packagePos := set.Position(file.Package); packagePos.IsValid() {
		packageLine = packagePos.Line
	}
	// Load comment groups and map them to their ending line number.
	// Assume a comment belongs to a command line flag declaration if it
	// ends on the previous line of the flag declaration.
	comments := make(map[int]*ast.CommentGroup)
	for _, group := range file.Comments {
		pos := set.Position(group.Pos())
		end := set.Position(group.End())
		if pos.Line < packageLine {
			f.Doc = group.Text()
		}
		comments[end.Line] = group
	}
	// Memorize options with variable names
	bound := make(map[string]*Option)
	// Load all options in source file. This means, detect and parse
	// all flag.Bool, flag.Duration, ... calls.
	ast.Inspect(file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if opt, err := NewOptionFromCallExpr(set, call); err == nil {
				// Check if we have a comment that belongs to option
				if comment, ok := comments[opt.Line-1]; ok {
					opt.Doc = comment.Text()
				}
				// Check if we already encountered an option bound to the
				// variable.
				if opt.Variable != "" {
					if reg, ok := bound[opt.Variable]; ok {
						// Merge currrent option with the one we already found
						reg.merge(opt)
						// Don't add the current option to the list, since the list
						// already contains the struct stored in the map.
						return true
					} else {
						// Register variable and the proceed to add option
						// struct to the options list
						bound[opt.Variable] = opt
					}
				}
				f.Options = append(f.Options, opt)
			}
		}
		return true
	})
	return nil
}
