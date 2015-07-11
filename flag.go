package main

import (
	"errors"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

var (
	ErrNotSupported = errors.New("not supported")
	ErrParse        = errors.New("parser error")
)

var (
	regexFlagCall  = regexp.MustCompile(`^[Ff]lags?\.(Bool|Duration|Float|Float64|Int|Int64|String|Uint|Uint64)(Var)?$`)
	regexBackquote = regexp.MustCompile("`(.*)`")
)

type FunctionCall struct {
	Line int
	Expr string
	Name string
	Args []string
}

func NewFunctionCall(fs *token.FileSet, n *ast.CallExpr) (*FunctionCall, error) {
	parts := make([]string, 0)
	ast.Inspect(n.Fun, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch obj := n.(type) {
		case *ast.SelectorExpr:
			// do nothing, just avoid the default case
			break
		case *ast.Ident:
			parts = append(parts, obj.Name)
		default:
			return false
		}
		return true
	})
	if len(parts) == 0 {
		return nil, ErrNotSupported
	}
	args := make([]string, len(n.Args))
	for i, arg := range n.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.Ident:
				args[i] = t.Name
			case *ast.BasicLit:
				args[i] = t.Value
			}
			return true
		})
	}
	return &FunctionCall{
		Line: fs.Position(n.Pos()).Line,
		Expr: strings.Join(parts, "."),
		Name: parts[len(parts)-1],
		Args: args,
	}, nil
}

type FlagType uint32

const (
	UnkownFlag FlagType = iota
	BoolFlag
	DurationFlag
	FloatFlag
	IntFlag
	UintFlag
	StringFlag
)

var flagTypes = map[string]FlagType{
	"Bool":     BoolFlag,
	"Duration": DurationFlag,
	"Float":    FloatFlag,
	"Float64":  FloatFlag,
	"Int":      IntFlag,
	"Int64":    IntFlag,
	"String":   StringFlag,
	"Uint":     UintFlag,
	"Uint64":   UintFlag,
}

// flagParam maps flag types to their default parameter names. These
// parameter names will be used if the usage string does not contain a
// backquoted name.
var flagParam = map[FlagType]string{
	DurationFlag: "duration",
	FloatFlag:    "float",
	IntFlag:      "int",
	StringFlag:   "string",
	UintFlag:     "uint",
}

type Flag struct {
	Type     FlagType
	Line     int
	Variable string // Pointer name (only set for ...Var() calls)
	Name     string // Name of the flag
	Short    string // Shorthand name of the flag
	Usage    string // Usage of the flag
	Param    string // User specified parameter name (back-quoted word in usage)
	Doc      string // Comment above flag declaration
}

const Quotes = "\"`"

func NewFlag(fs *token.FileSet, n *ast.CallExpr) (*Flag, error) {
	call, err := NewFunctionCall(fs, n)
	if err != nil {
		return nil, err
	}
	match := regexFlagCall.FindStringSubmatch(call.Expr)
	if match == nil {
		return nil, ErrNotSupported
	}
	// Pad to 4 arguments.
	if len(call.Args) == 3 {
		call.Args = append([]string{""}, call.Args...)
	}
	result := &Flag{
		Type:     flagTypes[match[1]],
		Line:     call.Line,
		Variable: call.Args[0],
		Name:     strings.Trim(call.Args[1], Quotes),
		Usage:    strings.Trim(call.Args[3], Quotes),
	}
	// Check if there's a backquoted parameter name in the usage string.
	if match := regexBackquote.FindStringSubmatch(result.Usage); match != nil {
		result.Param = match[1]
		result.Usage = strings.Replace(result.Usage, "`", "", -1)
	} else {
		result.Param = flagParam[result.Type]
	}
	return result, nil
}

func assignIfEmpty(d *string, v string) {
	if *d == "" {
		*d = v
	}
}

func (o *Flag) merge(v *Flag) {
	if len(o.Name) < len(v.Name) {
		o.Short = o.Name
		o.Name = v.Name
	} else {
		o.Short = v.Name
	}
	assignIfEmpty(&o.Doc, v.Doc)
	assignIfEmpty(&o.Param, v.Param)
	assignIfEmpty(&o.Usage, v.Usage)
}
