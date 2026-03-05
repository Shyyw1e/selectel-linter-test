package rules

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type LogCall struct {
	Kind string
	Level string
	Msg string 
	MsgExpr ast.Expr
	Call *ast.CallExpr
}

type RuleContext struct {
	SensitiveKeywords []string
}

type Issue struct {
	RuleID string
	Message string
	Node ast.Node
	Fixes []analysis.SuggestedFix
}

type Rule interface {
	ID() string
	Check(pass *analysis.Pass, lc LogCall, ctx RuleContext) []Issue
}