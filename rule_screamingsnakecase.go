package main

import (
	"sync"
	"unicode"
)

type RuleScreamingSnakeCase struct {
	Name string
	*sync.RWMutex
}

func (rule *RuleScreamingSnakeCase) Init() Rule {
	rule.Name = "screamingsnakecase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RuleScreamingSnakeCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *RuleScreamingSnakeCase) SetParameters(params []string) error {
	return nil
}

// Validate checks if string is screaming sneak case
// false if rune is no uppercase letter, digit or _
func (rule *RuleScreamingSnakeCase) Validate(value string) (bool, error) {
	for _, c := range value {
		if c == 95 || unicode.IsDigit(c) { // 95 => _
			continue
		}

		if !unicode.IsLetter(c) {
			return false, nil
		}

		if !unicode.IsUpper(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *RuleScreamingSnakeCase) GetErrorMessage() string {
	return rule.GetName()
}
