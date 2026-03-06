package analyser

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/Shyyw1e/selectel-linter-test/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var supportedLevels = map[string]struct{}{
	"debug": {},
	"info":  {},
	"warn":  {},
	"error": {},
}

func ParseLogCall(pass *analysis.Pass, call *ast.CallExpr) (rules.LogCall, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return rules.LogCall{}, false
	}

	level := strings.ToLower(sel.Sel.Name)
	if _, ok := supportedLevels[level]; !ok {
		return rules.LogCall{}, false
	}
	if len(call.Args) == 0 {
		return rules.LogCall{}, false
	}

	kind, ok := detectLoggerKind(pass, sel.X)
	if !ok {
		return rules.LogCall{}, false
	}

	msgExpr := call.Args[0]
	msg, ok := ExtractString(msgExpr)
	if !ok {
		return rules.LogCall{}, false
	}

	return rules.LogCall{
		Kind:    kind,
		Level:   level,
		Msg:     msg,
		MsgExpr: msgExpr,
		Call:    call,
	}, true
}

func detectLoggerKind(pass *analysis.Pass, recv ast.Expr) (string, bool) {
	if pass == nil || pass.TypesInfo == nil {
		if id, ok := recv.(*ast.Ident); ok && id.Name == "slog" {
			return "slog", true
		}
		return "zap", true
	}

	if id, ok := recv.(*ast.Ident); ok {
		if obj, ok := pass.TypesInfo.Uses[id].(*types.PkgName); ok {
			if obj.Imported() != nil && obj.Imported().Path() == "log/slog" {
				return "slog", true
			}
		}
	}

	tv, ok := pass.TypesInfo.Types[recv]
	if !ok || tv.Type == nil {
		return "", false
	}
	return loggerKindFromType(tv.Type)
}

func loggerKindFromType(t types.Type) (string, bool) {
	for {
		p, ok := t.(*types.Pointer)
		if !ok {
			break
		}
		t = p.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return "", false
	}

	path := named.Obj().Pkg().Path()
	name := named.Obj().Name()

	if path == "log/slog" && name == "Logger" {
		return "slog", true
	}
	if path == "go.uber.org/zap" && (name == "Logger" || name == "SugaredLogger") {
		return "zap", true
	}

	return "", false
}

func ExtractString(expr ast.Expr) (string, bool) {
	type frame struct {
		expr ast.Expr
	}

	stack := []frame{{expr: expr}}
	parts := make([]string, 0, 4)

	for len(stack) > 0 {
		n := len(stack) - 1
		cur := stack[n]
		stack = stack[:n]

		switch v := cur.expr.(type) {
		case *ast.ParenExpr:
			stack = append(stack, frame{expr: v.X})

		case *ast.BasicLit:
			if v.Kind != token.STRING {
				return "", false
			}
			s, err := strconv.Unquote(v.Value)
			if err != nil {
				return "", false
			}
			parts = append(parts, s)

		case *ast.BinaryExpr:
			if v.Op != token.ADD {
				return "", false
			}
			stack = append(stack, frame{expr: v.Y})
			stack = append(stack, frame{expr: v.X})

		default:
			return "", false
		}
	}

	return strings.Join(parts, ""), true
}
