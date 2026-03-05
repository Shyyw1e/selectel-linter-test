package analyser

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestExtractString(t *testing.T) {
	tests := []struct {
		name    string
		expr    ast.Expr
		want    string
		wantOK  bool
	}{
		{
			name:   "basic literal",
			expr:   &ast.BasicLit{Kind: token.STRING, Value: `"hello"`},
			want:   "hello",
			wantOK: true,
		},
		{
			name: "concat literals",
			expr: &ast.BinaryExpr{
				Op: token.ADD,
				X:  &ast.BasicLit{Kind: token.STRING, Value: `"hel"`},
				Y:  &ast.BasicLit{Kind: token.STRING, Value: `"lo"`},
			},
			want:   "hello",
			wantOK: true,
		},
		{
			name: "non-string literal",
			expr: &ast.BasicLit{Kind: token.INT, Value: "42"},
			want: "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractString(tt.expr)
			if ok != tt.wantOK {
				t.Fatalf("ok mismatch: got %v want %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("value mismatch: got %q want %q", got, tt.want)
			}
		})
	}
}
