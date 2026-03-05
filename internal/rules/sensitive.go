package rules

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

type SensitiveRule struct{}

func (r SensitiveRule) ID() string {
	return RuleNoSensitive
}

func (r SensitiveRule) Check(_ *analysis.Pass, lc LogCall, ctx RuleContext) []Issue {
	if lc.Msg == "" || lc.MsgExpr == nil {
		return nil
	}

	msg := strings.ToLower(lc.Msg)
	for _, kw := range ctx.SensitiveKeywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw == "" {
			continue
		}
		if strings.Contains(msg, kw) {
			return []Issue{
				{	
					RuleID:  r.ID(),
					Message: "log message must not contain sensitive data",
					Node:    lc.MsgExpr,
				},
			}
		}
	}

	return nil
}
