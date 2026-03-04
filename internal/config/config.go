package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Shyyw1e/selectel-linter-test/internal/rules"
)


type Config struct {
	EnabledRules		map[string]bool
	SensitiveKeywords	[]string
	AllowPunctuation	bool
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
		AllowPunctuation: false,
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

	normalized := make([]string, 0, len(c.SensitiveKeywords))
	for _, kw := range c.SensitiveKeywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw == "" {
			continue
		}
		normalized = append(normalized, kw)
	}
	if len(normalized) == 0 {
		return errors.New("sensitive_keywords is empty after normalisation")
	}
	c.SensitiveKeywords = normalized

	return nil
}