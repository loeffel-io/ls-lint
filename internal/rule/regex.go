package rule

import (
	"fmt"
	"regexp"
	"sync"
)

type Regex struct {
	name         string
	exclusive    bool
	regexPattern string
	*sync.RWMutex
}

func (rule *Regex) Init() Rule {
	rule.name = "regex"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Regex) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

// 0 = regex pattern
func (rule *Regex) SetParameters(params []string) error {
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

func (rule *Regex) GetParameters() []string {
	return []string{rule.regexPattern}
}

func (rule *Regex) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if full string matches regex
func (rule *Regex) Validate(value string, fail bool) (bool, error) {
	return regexp.MatchString(fmt.Sprintf("^%s$", rule.getRegexPattern()), value)
}

func (rule *Regex) getRegexPattern() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.regexPattern
}

func (rule *Regex) GetErrorMessage() string {
	return fmt.Sprintf("%s:%s", rule.GetName(), rule.getRegexPattern())
}

func (rule *Regex) Copy() Rule {
	rule.RLock()
	defer rule.RUnlock()

	c := new(Regex)
	c.Init()
	c.regexPattern = rule.regexPattern

	return c
}
