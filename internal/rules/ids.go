package rules

const (
	RuleLowercaseStart = "lowercase-start"
	RuleEnglishOnly    = "english-only"
	RuleNoSpecialChars = "no-special-chars"
	RuleNoSensitive    = "no-sensitive-data"
)

func KnownIDs() map[string]struct{} { // структура, потому что под нее выделяем 0 байт
	return map[string]struct{}{
		RuleLowercaseStart: {},
		RuleEnglishOnly:    {},
		RuleNoSpecialChars: {},
		RuleNoSensitive:    {},
	}
}
