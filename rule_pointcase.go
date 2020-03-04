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

// Validate checks if string is "point case"
// false if rune is no lowercase letter or .
func (rule *RulePointCase) Validate(value string) (bool, error) {
	for _, c := range value {
		if c == 46 { // .
			continue
		}

		if unicode.IsLetter(c) && !unicode.IsLower(c) {
			return false, nil
		}
	}

	return true, nil
}
