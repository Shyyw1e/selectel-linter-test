package rules

import (
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type NoSpecialRule struct{}

func (r NoSpecialRule) ID() string {
	return RuleNoSpecialChars
}

func (r NoSpecialRule) Check(_ *analysis.Pass, lc LogCall, ctx RuleContext) []Issue {
	if lc.Msg == "" || lc.MsgExpr == nil {
		return nil
	}

	for _, ch := range lc.Msg {
		if isAllowedChar(ch) {
			continue
		}

		return []Issue{
			{
				RuleID:  r.ID(),
				Message: "log message must not contain special characters or emoji",
				Node:    lc.MsgExpr,
			},
		}
	}

	return nil
}

func isAllowedChar(ch rune) bool {
	if unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSpace(ch) {
		return true
	}
	return false
}
