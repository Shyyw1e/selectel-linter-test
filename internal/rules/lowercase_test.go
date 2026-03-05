package rules

import (
	"go/ast"
	"testing"
)

func TestLowercaseRule_Check(t *testing.T) {
	r := LowercaseRule{}
	node := ast.NewIdent("msg")

	tests := []struct {
		name    string
		msg     string
		wantHit bool
	}{
		{"ok lowercase", "starting server", false},
		{"bad uppercase", "Starting server", true},
		{"ok leading space then lowercase", "  starting server", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := r.Check(nil, LogCall{Msg: tt.msg, MsgExpr: node}, RuleContext{})
			if tt.wantHit && len(issues) == 0 {
				t.Fatalf("expected issue, got none")
			}
			if !tt.wantHit && len(issues) != 0 {
				t.Fatalf("expected no issues, got %d", len(issues))
			}
		})
	}
}
