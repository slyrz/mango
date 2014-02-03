package source

import (
	"errors"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

var (
	ErrNoSup = errors.New("not supported")
	ErrParse = errors.New("parser error")
)

var (
	reFlag = regexp.MustCompile("^[Ff]lags?\\.(Bool|Duration|Float|Float64|Int|Int64|String|Uint|Uint64)(Var)?$")
)

func getFunctionName(e ast.Expr) (string, error) {
	parts := make([]string, 0, 4)
	ast.Inspect(e, func(n ast.Node) bool {
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

	// We are only interested in Package.Name type of calls...
	if len(parts) != 2 {
		return "", ErrNoSup
	}
	return strings.Join(parts, "."), nil
}

func getArgValue(arg ast.Expr) string {
	result := ""
	ast.Inspect(arg, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			result = lit.Value
		}
		return false
	})
	return result
}

func getArgIdentifier(arg ast.Expr) string {
	result := ""
	ast.Inspect(arg, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			result = ident.Name
		}
		return true
	})
	return result
}

type Option struct {
	Line     int
	Type     string // Bool | Duration | Float ...
	Variable string // Pointer name (only for ...Var() calls)
	Name     string // Name of the flag
	Short    string // Shorthand name of the flag
	Usage    string // Usage of the flag
	Doc      string // Comment near flag
}

func NewOptionFromCallExpr(fs *token.FileSet, n *ast.CallExpr) (*Option, error) {
	funcName, err := getFunctionName(n.Fun)
	if err != nil {
		return nil, err
	}

	match := reFlag.FindStringSubmatch(funcName)
	if match == nil {
		return nil, ErrNoSup
	}

	option := new(Option)
	option.Type = match[1]

	argsLen := 3
	argsOff := 0
	if match[2] == "Var" {
		argsLen++
		argsOff++
	}

	if len(n.Args) != argsLen {
		return nil, ErrNoSup
	}

	if argsLen == 4 {
		option.Variable = getArgIdentifier(n.Args[0])
		if len(option.Variable) == 0 {
			return nil, ErrParse
		}
	}

	option.Name = getArgValue(n.Args[argsOff+0])
	option.Usage = getArgValue(n.Args[argsOff+2])

	if len(option.Name) == 0 || len(option.Usage) == 0 {
		return nil, ErrParse
	}

	option.Name = strings.Trim(option.Name, "\"`")
	option.Usage = strings.Trim(option.Usage, "\"`")

	option.Line = fs.Position(n.Pos()).Line
	return option, nil
}

func (o *Option) merge(v *Option) {
	if len(o.Name) < len(v.Name) {
		o.Short = o.Name
		o.Name = v.Name
	} else {
		o.Short = v.Name
	}

	if len(o.Doc) == 0 {
		o.Doc = v.Doc
	}

	if len(o.Usage) == 0 {
		o.Usage = v.Usage
	}
}
