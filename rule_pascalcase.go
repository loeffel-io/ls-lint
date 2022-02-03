package main

import (
	"regexp"
	"sync"
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
// false if rune is no letter and no digit
// false if first rune is not upper
func (rule *RulePascalCase) Validate(value string) (bool, error) {
	return regexp.MatchString("^([A-Z][a-z]*[0-9]*)+$", value)
}

func (rule *RulePascalCase) GetErrorMessage() string {
	return rule.GetName()
}
