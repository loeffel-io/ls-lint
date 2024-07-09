package rule

import (
	"sync"
	"unicode"
)

type PascalCaseAcronym struct {
	Name string
	*sync.RWMutex
}

func (rule *PascalCaseAcronym) Init() Rule {
	rule.Name = "pascalcaseacronym"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *PascalCaseAcronym) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *PascalCaseAcronym) SetParameters([]string) error {
	return nil
}

func (rule *PascalCaseAcronym) GetParameters() []string {
	return nil
}

// Validate checks if string is pascal case
// false if rune is no letter and no digit
// false if first rune is not upper
// allows up to five consecutive upper letters
// allow cases like CTOClown.go or NASAImages/
func (rule *PascalCaseAcronym) Validate(value string) (bool, error) {
	upperStreak := 0

	for i, c := range value {
		// must be letter or digit
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false, nil
		}

		// first rune must be upper
		if i == 0 && unicode.IsLower(c) {
			return false, nil
		}

		if unicode.IsUpper(c) {
			upperStreak++
			if i == 0 {
				continue
			}

			if upperStreak > 5 {
				return false, nil
			}
		} else {
			upperStreak = 0
		}
	}

	return true, nil
}

func (rule *PascalCaseAcronym) GetErrorMessage() string {
	return rule.GetName()
}
