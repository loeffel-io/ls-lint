package rule

import (
	"sync"
	"unicode"
)

type CamelCase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *CamelCase) Init() Rule {
	rule.name = "camelcase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *CamelCase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *CamelCase) SetParameters(params []string) error {
	return nil
}

func (rule *CamelCase) GetParameters() []string {
	return nil
}

func (rule *CamelCase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if string is camel case
// false if rune is no letter and no digit
func (rule *CamelCase) Validate(value string, fail bool) (bool, error) {
	for i, c := range value {
		// must be letter or digit
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false, nil
		}

		if unicode.IsUpper(c) {
			// first rune cannot be upper
			if i == 0 {
				return false, nil
			}

			// rune -1 can be digit
			if unicode.IsDigit(rune(value[i-1])) {
				continue
			}

			// allow cases like ssrVFor.ts
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

func (rule *CamelCase) GetErrorMessage() string {
	return rule.GetName()
}
