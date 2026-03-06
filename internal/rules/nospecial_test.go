package rules

import (
	"go/ast"
	"go/token"
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

func TestLowercaseRule_FixForStringLiteral(t *testing.T) {
	r := LowercaseRule{}
	lit := &ast.BasicLit{Kind: token.STRING, Value: `"Starting server"`}

	issues := r.Check(nil, LogCall{
		Msg:     "Starting server",
		MsgExpr: lit,
	}, RuleContext{})

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if len(issues[0].Fixes) != 1 {
		t.Fatalf("expected 1 suggested fix, got %d", len(issues[0].Fixes))
	}
	edits := issues[0].Fixes[0].TextEdits
	if len(edits) != 1 {
		t.Fatalf("expected 1 text edit, got %d", len(edits))
	}
	if string(edits[0].NewText) != `"starting server"` {
		t.Fatalf("unexpected fix text: %s", string(edits[0].NewText))
	}
}
