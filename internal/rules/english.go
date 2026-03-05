package rules

import (
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type EnglishRule struct{}

func (r EnglishRule) ID() string {
	return RuleEnglishOnly
}

func (r EnglishRule) Check(_ *analysis.Pass, lc LogCall, _ RuleContext) []Issue {
	if lc.Msg == "" || lc.MsgExpr == nil {
		return nil
	}

	for _, ch := range lc.Msg {
		if unicode.In(ch, unicode.Cyrillic) {
			return []Issue{
				{
					RuleID:  r.ID(),
					Message: "log message must be in English",
					Node:    lc.MsgExpr,
				},
			}
		}
	}

	return nil
}
