package rule

import (
	"fmt"
	"regexp"
	"sync"
)

type Regex struct {
	Name         string
	RegexPattern string
	*sync.RWMutex
}

func (rule *Regex) Init() Rule {
	rule.Name = "regex"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Regex) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
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

	rule.RegexPattern = params[0]
	return nil
}

func (rule *Regex) GetParameters() []string {
	return []string{rule.RegexPattern}
}

// Validate checks if full string matches regex
func (rule *Regex) Validate(value string) (bool, error) {
	return regexp.MatchString(fmt.Sprintf("^%s$", rule.getRegexPattern()), value)
}

func (rule *Regex) getRegexPattern() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.RegexPattern
}

func (rule *Regex) GetErrorMessage() string {
	return fmt.Sprintf("%s (%s)", rule.GetName(), rule.getRegexPattern())
}
