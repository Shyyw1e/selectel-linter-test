package rules

import (
	"unicode"

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

	for _, ch := range lc.Msg {
		if unicode.IsSpace(ch) {
			continue
		}
		if !unicode.IsLower(ch) {
			return []Issue{
				{
					RuleID:  r.ID(),
					Message: "log message must start with a lowercase letter",
					Node:    lc.MsgExpr,
				},
			}
		}
		return nil
	}

	return nil
}