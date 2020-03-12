package main

import (
	"sync"
	"unicode"
)

type RulePascalCase struct {
	Name string
	*sync.RWMutex
}

func (rule *RulePascalCase) Init() Rule {
	rule.Name = "pascalcase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RulePascalCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *RulePascalCase) SetParameters(params []string) error {
	return nil
}

// Validate checks if string is pascal case
// false if rune is no letter
// false if first rune is not upper
func (rule *RulePascalCase) Validate(value string) (bool, error) {
	for i, c := range value {
		if !unicode.IsLetter(c) {
			return false, nil
		}

		if i == 0 && unicode.IsLower(c) {
			return false, nil
		}

		if unicode.IsUpper(c) {
			if i == 0 {
				continue
			}

			if !unicode.IsLower(rune(value[i-1])) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (rule *RulePascalCase) GetErrorMessage() string {
	return rule.GetName()
}
