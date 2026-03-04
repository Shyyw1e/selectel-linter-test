package config

import "testing"

func TestDefaultConfigIsValid(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestValidateFailsOnUnknownRule(t *testing.T) {
	cfg := Default()
	cfg.EnabledRules["unknown-rule"] = true

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unknown rule")
	}
}

func TestValidateFailsOnEmptyEnabledRules(t *testing.T) {
	cfg := Default()
	cfg.EnabledRules = map[string]bool{}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty enabled_rules")
	}
}

func TestValidateFailsOnEmptySensitiveKeywordsAfterNormalization(t *testing.T) {
	cfg := Default()
	cfg.SensitiveKeywords = []string{" ", "\t", ""}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty sensitive_keywords after normalization")
	}
}
