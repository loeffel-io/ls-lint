package rule

import (
	"sync"
	"unicode"
)

type UppercaseDigit struct {
	Name string
	*sync.RWMutex
}

func (rule *UppercaseDigit) Init() Rule {
	rule.Name = "uppercasedigit"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *UppercaseDigit) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *UppercaseDigit) SetParameters([]string) error {
	return nil
}

func (rule *UppercaseDigit) GetParameters() []string {
	return nil
}

// Validate checks if string is all uppercase or digit
// false if rune is no uppercase letter or digit
func (rule *UppercaseDigit) Validate(value string) (bool, error) {
	for _, c := range value {
		if unicode.IsDigit(c) {
			continue
		}
		if !unicode.IsLetter(c) || !unicode.IsUpper(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *UppercaseDigit) GetErrorMessage() string {
	return rule.GetName()
}
