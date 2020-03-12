package main

import (
	"sync"
	"unicode"
)

type RuleLowercase struct {
	Name string
	*sync.RWMutex
}

func (rule *RuleLowercase) Init() Rule {
	rule.Name = "lowercase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RuleLowercase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *RuleLowercase) SetParameters(params []string) error {
	return nil
}

// Validate checks if every letter is lower
func (rule *RuleLowercase) Validate(value string) (bool, error) {
	for _, c := range value {
		if unicode.IsLetter(c) && !unicode.IsLower(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *RuleLowercase) GetErrorMessage() string {
	return rule.GetName()
}
