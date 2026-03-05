package analyser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestParseLogCall(t *testing.T) {
	src := `
package p

import "log/slog"

type logger struct{}
func (logger) Info(msg string, fields ...any) {}
func (logger) Error(msg string, fields ...any) {}
func (logger) Printf(format string, args ...any) {}

func f() {
	slog.Info("starting server")
	slog.Error("failed " + "connect")

	var z logger
	z.Info("zap style info")
	z.Error("zap style error")

	z.Printf("not a level")
	slog.Info()
	slog.Info(42)
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatalf("parse file: %v", err)
	}

	var calls []*ast.CallExpr
	ast.Inspect(file, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			calls = append(calls, c)
		}
		return true
	})

	if len(calls) == 0 {
		t.Fatal("no call expressions found")
	}

	var got []string
	for _, c := range calls {
		lc, ok := ParseLogCall(nil, c)
		if !ok {
			continue
		}
		got = append(got, lc.Kind+":"+lc.Level+":"+lc.Msg)
	}

	want := map[string]bool{
		"slog:info:starting server": true,
		"slog:error:failed connect": true,
		"zap:info:zap style info":   true,
		"zap:error:zap style error": true,
	}

	if len(got) != len(want) {
		t.Fatalf("unexpected parsed call count: got=%d want=%d values=%v", len(got), len(want), got)
	}

	for _, v := range got {
		if !want[v] {
			t.Fatalf("unexpected parsed call: %s", v)
		}
	}
}
