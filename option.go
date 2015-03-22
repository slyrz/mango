package main

import (
	"errors"
	"fmt"
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

type OptionType uint32

const (
	UnkownOption OptionType = iota
	BoolOption
	DurationOption
	FloatOption
	IntOption
	UintOption
	StringOption
)

var optionTypes = map[string]OptionType{
	"Bool":     BoolOption,
	"Duration": DurationOption,
	"Float":    FloatOption,
	"Float64":  FloatOption,
	"Int":      IntOption,
	"Int64":    IntOption,
	"String":   StringOption,
	"Uint":     UintOption,
	"Uint64":   UintOption,
}

// optionParam maps option types to their default parameter names. These
// parameter names will be used if the usage string does not contain a
// backquoted name.
var optionParam = map[OptionType]string{
	DurationOption: "duration",
	FloatOption:    "float",
	IntOption:      "int",
	StringOption:   "string",
	UintOption:     "uint",
}

type Option struct {
	Type     OptionType
	Line     int
	Variable string // Pointer name (only set for ...Var() calls)
	Name     string // Name of the flag
	Short    string // Shorthand name of the flag
	Usage    string // Usage of the flag
	Param    string // User specified parameter name (back-quoted word in usage)
	Doc      string // Comment above flag declaration
}

const Quotes = "\"`"

func NewOptionFromCallExpr(fs *token.FileSet, n *ast.CallExpr) (*Option, error) {
	call, err := NewFunctionCall(fs, n)
	if err != nil {
		fmt.Println("ERROR")
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
	result := &Option{
		Type:     optionTypes[match[1]],
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
		result.Param = optionParam[result.Type]
	}
	return result, nil
}

func assignIfEmpty(d *string, v string) {
	if *d == "" {
		*d = v
	}
}

func (o *Option) merge(v *Option) {
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
