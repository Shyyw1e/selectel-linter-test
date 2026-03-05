package rules

func DefaultRules() []Rule {
	return []Rule{
		LowercaseRule{},
		EnglishRule{},
		NoSpecialRule{},
		SensitiveRule{},
	}
}
