package rule

import (
	"sync"
	"unicode"
)

type PascalCase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *PascalCase) Init() Rule {
	rule.name = "pascalcase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *PascalCase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *PascalCase) SetParameters(params []string) error {
	return nil
}

func (rule *PascalCase) GetParameters() []string {
	return nil
}

func (rule *PascalCase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if string is pascal case
// false if rune is no letter and no digit
// false if first rune is not upper
func (rule *PascalCase) Validate(str string, _ string, _ bool) (bool, error) {
	// we will iterate over the string as runes
	// it allows us to get the previous rune easily
	value := []rune(str)
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
			if unicode.IsDigit(value[i-1]) {
				continue
			}

			// allow cases like SsrVFor.ts
			if i >= 2 && unicode.IsUpper(value[i-1]) && unicode.IsLower(value[i-2]) {
				continue
			}

			// rune -1 must be lower
			if !unicode.IsLower(value[i-1]) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (rule *PascalCase) GetErrorMessage() string {
	return rule.GetName()
}

func (rule *PascalCase) Copy() Rule {
	return rule
}
