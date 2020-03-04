package main

import (
	"sync"
	"unicode"
)

type RuleKebabCase struct {
	Name string
	*sync.RWMutex
}

func (rule *RuleKebabCase) Init() Rule {
	rule.Name = "kebabcase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RuleKebabCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

// Validate checks if string is kebab case
// false if rune is no lowercase letter or -
func (rule *RuleKebabCase) Validate(value string) (bool, error) {
	for _, c := range value {
		if c == 45 { // -
			continue
		}

		if !unicode.IsLetter(c) {
			return false, nil
		}

		if !unicode.IsLower(c) {
			return false, nil
		}
	}

	return true, nil
}
