package rule

import (
	"fmt"
	"regexp"
	"sync"
)

type NotRegex struct {
	name         string
	exclusive    bool
	regexPattern string
	*sync.RWMutex
}

func (rule *NotRegex) Init() Rule {
	rule.name = "not_regex"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *NotRegex) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

// 0 = regex pattern
func (rule *NotRegex) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	if len(params) == 0 {
		return fmt.Errorf("regex pattern not exists")
	}

	if params[0] == "" {
		return fmt.Errorf("regex pattern is empty")
	}

	rule.regexPattern = params[0]
	return nil
}

func (rule *NotRegex) GetParameters() []string {
	return []string{rule.regexPattern}
}

func (rule *NotRegex) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if the full string does NOT match the regex
func (rule *NotRegex) Validate(value string, fail bool) (bool, error) {
	match, err := regexp.MatchString(fmt.Sprintf("^%s$", rule.getRegexPattern()), value)
	return !match, err
}

func (rule *NotRegex) getRegexPattern() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.regexPattern
}

func (rule *NotRegex) GetErrorMessage() string {
	return fmt.Sprintf("%s:%s", rule.GetName(), rule.getRegexPattern())
}

func (rule *NotRegex) Copy() Rule {
	rule.RLock()
	defer rule.RUnlock()

	c := new(Regex)
	c.Init()
	c.regexPattern = rule.regexPattern

	return c
}
