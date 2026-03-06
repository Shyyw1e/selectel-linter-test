package rules

import (
	"go/ast"
	"testing"
)

func TestSensitiveRule_Check(t *testing.T) {
	r := SensitiveRule{}
	node := ast.NewIdent("msg")

	ctx := RuleContext{
		SensitiveKeywords: []string{"password", "token", "api_key"},
		SensitivePatterns: []string{`\bapi[-_]?key\b`, `\bbearer\b`},
	}

	tests := []struct {
		name    string
		msg     string
		wantHit bool
	}{
		{"ok safe", "user authenticated successfully", false},
		{"bad keyword password", "user password is invalid", true},
		{"bad keyword mixed case", "Token validated", true},
		{"bad pattern api-key", "api-key rotated", true},
		{"bad pattern bearer", "authorization bearer abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := r.Check(nil, LogCall{Msg: tt.msg, MsgExpr: node}, ctx)
			if tt.wantHit && len(issues) == 0 {
				t.Fatalf("expected issue, got none")
			}
			if !tt.wantHit && len(issues) != 0 {
				t.Fatalf("expected no issues, got %d", len(issues))
			}
		})
	}
}
