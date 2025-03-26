package rule

import (
	"sync"
	"unicode"
)

type Lowercase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *Lowercase) Init() Rule {
	rule.name = "lowercase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Lowercase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *Lowercase) SetParameters(params []string) error {
	return nil
}

func (rule *Lowercase) GetParameters() []string {
	return nil
}

func (rule *Lowercase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if every letter is lower
func (rule *Lowercase) Validate(value string, _ string, _ bool) (bool, error) {
	for _, c := range value {
		if unicode.IsLetter(c) && !unicode.IsLower(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *Lowercase) GetErrorMessage() string {
	return rule.GetName()
}

func (rule *Lowercase) Copy() Rule {
	return rule
}
