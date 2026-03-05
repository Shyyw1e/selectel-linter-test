package analyser

import (
	"go/ast"

	"github.com/Shyyw1e/selectel-linter-test/internal/config"
	"golang.org/x/tools/go/analysis"
)

type Analyser struct {
	cfg *config.Config
}

func NewAnalyser(cfg *config.Config) *analysis.Analyzer {
	if cfg == nil {
		def := config.Default()
		cfg = def
	}

	a := &Analyser{cfg: cfg}

	return &analysis.Analyzer{
		Name: "loglint",
		Doc: "checls log message style rules for slog/zap",
		Run: a.run,
	}
}

func (a *Analyser) run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			_ = call
			return true
		})
	}

	return nil, nil
}