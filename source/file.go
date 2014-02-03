package source

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
)

var (
	ErrFileType = errors.New("not a Go file")
)

type File struct {
	Options []*Option // Options found in file.
	Name    string    // Name of file.
	Doc     string    // Comment at top of file.

	// Unexported fields
	fileSet  *token.FileSet
	file     *ast.File
	comments map[int]*ast.CommentGroup // map ending line -> comment group
}

func NewFile(filePath string) (*File, error) {
	result := new(File)

	_, name := path.Split(filePath)

	// Check if path seems to be a .go file
	if ext := path.Ext(name); ext != ".go" {
		return nil, ErrFileType
	} else {
		result.Name = name[:len(name)-len(ext)] // remove ".go" to get name
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	result.fileSet = fileSet
	result.file = file

	// Load comment groups and map them to their ending line number.
	// We assume a comment belongs to a command line flag declaration if it
	// ends on the previous line of the flag declaration.
	result.comments = make(map[int]*ast.CommentGroup)
	for _, group := range file.Comments {
		pos := fileSet.Position(group.Pos())
		end := fileSet.Position(group.End())

		if pos.Line == 1 {
			result.Doc = group.Text()
		}
		result.comments[end.Line] = group
	}

	result.parseOptions()

	return result, nil
}

func (f *File) parseOptions() {
	options := make([]*Option, 0)

	// Memorize options with variable names
	bound := make(map[string]*Option)

	// Load all options in source file. This means, detect and parse
	// all flag.Bool, flag.Duration, ... calls.
	ast.Inspect(f.file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if opt, err := NewOptionFromCallExpr(f.fileSet, call); err == nil {
				// Check if we have a comment that belongs to option
				if comment, ok := f.comments[opt.Line-1]; ok {
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
				options = append(options, opt)
			}
		}
		return true
	})

	f.Options = options
}
