package plugin

import (
	"github.com/Shyyw1e/selectel-linter-test/internal/analyser"
	"github.com/Shyyw1e/selectel-linter-test/internal/config"
	"golang.org/x/tools/go/analysis"
)

func New(cfg *config.Config) *analysis.Analyzer {
	if cfg == nil {
		cfg = config.Default()
	}
	return analyser.NewAnalyser(cfg)
}
