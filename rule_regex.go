package main

import (
	"fmt"
	"regexp"
	"sync"
)

type RuleRegex struct {
	Name         string
	RegexPattern string
	*sync.RWMutex
}

func (rule *RuleRegex) Init() Rule {
	rule.Name = "regex"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RuleRegex) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

// 0 = regex pattern
func (rule *RuleRegex) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	if len(params) == 0 {
		return fmt.Errorf("regex pattern not exists")
	}

	if params[0] == "" {
		return fmt.Errorf("regex pattern is empty")
	}

	rule.RegexPattern = params[0]
	return nil
}

// Validate checks if string matches regex
func (rule *RuleRegex) Validate(value string) (bool, error) {
	return regexp.MatchString(rule.getRegexPattern(), value)
}

func (rule *RuleRegex) getRegexPattern() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.RegexPattern
}

func (rule *RuleRegex) GetErrorMessage() string {
	return fmt.Sprintf("%s (%s)", rule.GetName(), rule.getRegexPattern())
}
