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

// Validate checks if string is camel case
// false if rune is no letter
func (rule *RuleCamelCase) Validate(value string) (bool, error) {
	for i, c := range value {
		if !unicode.IsLetter(c) {
			return false, nil
		}

		if unicode.IsUpper(c) {
			if i == 0 || !unicode.IsLower(rune(value[i-1])) {
				return false, nil
			}
		}
	}

	return true, nil
}
