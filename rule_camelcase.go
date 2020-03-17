package main

import (
	"sync"
	"unicode"
)

type RuleCamelCase struct {
	Name string
	*sync.RWMutex
}

func (rule *RuleCamelCase) Init() Rule {
	rule.Name = "camelcase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RuleCamelCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *RuleCamelCase) SetParameters(params []string) error {
	return nil
}

// Validate checks if string is camel case
// false if rune is no letter and no digit
func (rule *RuleCamelCase) Validate(value string) (bool, error) {
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
			if unicode.IsUpper(rune(value[i-1])) && unicode.IsLower(rune(value[i-2])) {
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

func (rule *RuleCamelCase) GetErrorMessage() string {
	return rule.GetName()
}
