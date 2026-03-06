package rules

import (
	"regexp"
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

	msgLower := strings.ToLower(lc.Msg)

	for _, kw := range ctx.SensitiveKeywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw == "" {
			continue
		}
		if strings.Contains(msgLower, kw) {
			return []Issue{{
				RuleID:  r.ID(),
				Message: "log message must not contain sensitive data",
				Node:    lc.MsgExpr,
			}}
		}
	}

	for _, p := range ctx.SensitivePatterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		re, err := regexp.Compile(p)
		if err != nil {
			// config.Validate already checks this; ignore invalid pattern defensively.
			continue
		}
		if re.MatchString(msgLower) {
			return []Issue{{
				RuleID:  r.ID(),
				Message: "log message must not contain sensitive data",
				Node:    lc.MsgExpr,
			}}
		}
	}

	return nil
}
