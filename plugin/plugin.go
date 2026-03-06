package plugin

import (
	"github.com/Shyyw1e/selectel-linter-test/internal/analyser"
	"github.com/Shyyw1e/selectel-linter-test/internal/config"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type Settings struct {
	EnabledRules      map[string]bool `json:"enabled_rules" mapstructure:"enabled_rules"`
	SensitiveKeywords []string        `json:"sensitive_keywords" mapstructure:"sensitive_keywords"`
	SensitivePatterns []string        `json:"sensitive_patterns" mapstructure:"sensitive_patterns"`
}

type Plugin struct {
	cfg *config.Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	cfg := config.Default()

	if len(s.EnabledRules) > 0 {
		cfg.EnabledRules = s.EnabledRules
	}
	if len(s.SensitiveKeywords) > 0 {
		cfg.SensitiveKeywords = s.SensitiveKeywords
	}
	if len(s.SensitivePatterns) > 0 {
		cfg.SensitivePatterns = s.SensitivePatterns
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Plugin{cfg: cfg}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyser.NewAnalyser(p.cfg),
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
