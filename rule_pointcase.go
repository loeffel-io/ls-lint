package main

import (
	"sync"
	"unicode"
)

type RulePointCase struct {
	Name string
	*sync.RWMutex
}

func (rule *RulePointCase) Init() Rule {
	rule.Name = "pointcase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RulePointCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *RulePointCase) SetParameters(params []string) error {
	return nil
}

// Validate checks if string is "point case"
// false if rune is no lowercase letter, digit or .
func (rule *RulePointCase) Validate(value string) (bool, error) {
	for _, c := range value {
		if c == 46 || unicode.IsDigit(c) { // 46 => .
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

func (rule *RulePointCase) GetErrorMessage() string {
	return rule.GetName()
}
