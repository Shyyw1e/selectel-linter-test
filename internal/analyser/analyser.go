package analyser

import (
	"go/ast"

	"github.com/Shyyw1e/selectel-linter-test/internal/config"
	"github.com/Shyyw1e/selectel-linter-test/internal/rules"
	"golang.org/x/tools/go/analysis"
)

type Analyser struct {
	cfg      *config.Config
	ruleList []rules.Rule
}

func NewAnalyser(cfg *config.Config) *analysis.Analyzer {
	if cfg == nil {
		cfg = config.Default()
	}

	a := &Analyser{
		cfg:      cfg,
		ruleList: rules.DefaultRules(),
	}

	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "checks log message style rules for slog/zap",
		Run:  a.run,
	}
}

func (a *Analyser) run(pass *analysis.Pass) (any, error) {
	if err := a.cfg.Validate(); err != nil {
		return nil, err
	}

	ctx := rules.RuleContext{
		SensitiveKeywords: a.cfg.SensitiveKeywords,
		SensitivePatterns: a.cfg.SensitivePatterns,
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			lc, ok := ParseLogCall(pass, call)
			if !ok {
				return true
			}

			for _, rule := range a.ruleList {
				if enabled, exists := a.cfg.EnabledRules[rule.ID()]; exists && !enabled {
					continue
				}

				issues := rule.Check(pass, lc, ctx)
				for _, issue := range issues {
					if issue.Node == nil {
						continue
					}

					pass.Report(analysis.Diagnostic{
						Pos:            issue.Node.Pos(),
						End:            issue.Node.End(),
						Message:        issue.Message,
						SuggestedFixes: issue.Fixes,
					})
				}
			}

			return true
		})
	}

	return nil, nil
}
