package rules

import (
	"go/ast"
	"testing"
)

func TestEnglishRule_Check(t *testing.T) {
	r := EnglishRule{}
	node := ast.NewIdent("msg")

	tests := []struct {
		name    string
		msg     string
		wantHit bool
	}{
		{"ok english", "starting server", false},
		{"bad cyrillic", "запуск сервера", true},
		{"bad mixed", "start сервера", true},
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
