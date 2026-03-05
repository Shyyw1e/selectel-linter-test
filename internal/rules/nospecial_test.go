package rules

import (
	"go/ast"
	"testing"
)

func TestNoSpecialRule_Check(t *testing.T) {
	r := NoSpecialRule{}
	node := ast.NewIdent("msg")

	tests := []struct {
		name    string
		msg     string
		wantHit bool
	}{
		{"ok letters digits spaces", "server started 200", false},
		{"bad punctuation", "server started!", true},
		{"bad emoji", "server started 🚀", true},
		{"bad dots", "something went wrong...", true},
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
