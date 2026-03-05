package analyser

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/Shyyw1e/selectel-linter-test/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var supportedLevels = map[string]struct{} {
	"debug": {},
	"info": {},
	"warn": {},
	"error": {},
}

func ParseLogCall(_ *analysis.Pass, call *ast.CallExpr) (rules.LogCall, bool) {
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

	msgExpr := call.Args[0]
	msg, ok := ExtractString(msgExpr)
	if !ok {
		return rules.LogCall{}, false
	}

	kind := "unknown"
	if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "slog" {
		kind = "slog"
	} else {
		kind = "zap"
	}

	return rules.LogCall{
		Kind: kind,
		Level: level,
		Msg: msg,
		MsgExpr: msgExpr,
		Call: call,
	}, true
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
