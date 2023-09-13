package rule

import (
	"sync"
	"unicode"
)

type PascalCase struct {
	Name string
	*sync.RWMutex
}

func (rule *PascalCase) Init() Rule {
	rule.Name = "pascalcase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *PascalCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *PascalCase) SetParameters(params []string) error {
	return nil
}

func (rule *PascalCase) GetParameters() []string {
	return nil
}

// Validate checks if string is pascal case
// false if rune is no letter and no digit
// false if first rune is not upper
func (rule *PascalCase) Validate(value string) (bool, error) {
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
			if i == 0 {
				continue
			}

			// rune -1 can be digit
			if unicode.IsDigit(rune(value[i-1])) {
				continue
			}

			// allow cases like SsrVFor.ts
			if i >= 2 && unicode.IsUpper(rune(value[i-1])) && unicode.IsLower(rune(value[i-2])) {
				continue
			}

			// rune -1 must be lower
			if !unicode.IsLower(rune(value[i-1])) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (rule *PascalCase) GetErrorMessage() string {
	return rule.GetName()
}
