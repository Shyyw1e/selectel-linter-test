package rules

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

type LowercaseRule struct{}

func (r LowercaseRule) ID() string {
	return RuleLowercaseStart
}

func (r LowercaseRule) Check(_ *analysis.Pass, lc LogCall, _ RuleContext) []Issue {
	if lc.Msg == "" || lc.MsgExpr == nil {
		return nil
	}

	idx := firstNonSpaceIndex(lc.Msg)
	if idx < 0 {
		return nil
	}

	ch, size := utf8.DecodeRuneInString(lc.Msg[idx:])
	if ch == utf8.RuneError && size == 0 {
		return nil
	}
	if unicode.IsLower(ch) {
		return nil
	}

	issue := Issue{
		RuleID:  r.ID(),
		Message: "log message must start with a lowercase letter",
		Node:    lc.MsgExpr,
	}

	if lit, ok := lc.MsgExpr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		fixed := lc.Msg[:idx] + string(unicode.ToLower(ch)) + lc.Msg[idx+size:]
		quoted := strconv.Quote(fixed)

		issue.Fixes = []analysis.SuggestedFix{
			{
				Message: "make first letter lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     lit.Pos(),
						End:     lit.End(),
						NewText: []byte(quoted),
					},
				},
			},
		}
	}

	return []Issue{issue}
}

func firstNonSpaceIndex(s string) int {
	for i, ch := range s {
		if !unicode.IsSpace(ch) {
			return i
		}
	}
	return -1
}
