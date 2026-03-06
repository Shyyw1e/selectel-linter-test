package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Shyyw1e/selectel-linter-test/internal/rules"
)

type Config struct {
	EnabledRules      map[string]bool
	SensitiveKeywords []string
	SensitivePatterns []string
}

func Default() *Config {
	return &Config{
		EnabledRules: map[string]bool{
			rules.RuleLowercaseStart: true,
			rules.RuleEnglishOnly:    true,
			rules.RuleNoSpecialChars: true,
			rules.RuleNoSensitive:    true,
		},
		SensitiveKeywords: []string{
			"password",
			"passwd",
			"pwd",
			"token",
			"api_key",
			"apikey",
			"secret",
			"bearer",
			"credential",
		},
		SensitivePatterns: []string{
			`\bapi[_-]?key\b`,
			`\bbearer\b`,
			`\bpassword\b`,
		},
	}
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}
	if len(c.EnabledRules) == 0 {
		return errors.New("enabled_rules is empty")
	}

	known := rules.KnownIDs()
	for id := range c.EnabledRules {
		if _, ok := known[id]; !ok {
			return fmt.Errorf("unknown rule id: %s", id)
		}
	}

	normalizedKeywords := make([]string, 0, len(c.SensitiveKeywords))
	for _, kw := range c.SensitiveKeywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw == "" {
			continue
		}
		normalizedKeywords = append(normalizedKeywords, kw)
	}
	if len(normalizedKeywords) == 0 {
		return errors.New("sensitive_keywords is empty after normalization")
	}
	c.SensitiveKeywords = normalizedKeywords

	normalizedPatterns := make([]string, 0, len(c.SensitivePatterns))
	for _, p := range c.SensitivePatterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("invalid sensitive pattern %q: %w", p, err)
		}
		normalizedPatterns = append(normalizedPatterns, p)
	}
	c.SensitivePatterns = normalizedPatterns

	return nil
}
