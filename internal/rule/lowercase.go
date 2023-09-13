package rule

import (
	"sync"
	"unicode"
)

type Lowercase struct {
	Name string
	*sync.RWMutex
}

func (rule *Lowercase) Init() Rule {
	rule.Name = "lowercase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Lowercase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *Lowercase) SetParameters(params []string) error {
	return nil
}

func (rule *Lowercase) GetParameters() []string {
	return nil
}

// Validate checks if every letter is lower
func (rule *Lowercase) Validate(value string) (bool, error) {
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
