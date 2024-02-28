package rule

import (
	"fmt"
	"strings"
	"sync"
)

type Disallow struct {
	Name    string
	Message string
	*sync.RWMutex
}

func (rule *Disallow) Init() Rule {
	rule.Name = "disallow"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Disallow) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *Disallow) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	if len(params) > 0 && params[0] != "" {
		rule.Message = strings.TrimSpace(params[0])
	} else {
		rule.Message = ""
	}
	return nil
}

func (rule *Disallow) GetParameters() []string {
	return []string{rule.Message}
}

// Validate always returns false
// Any string that matches fails validation.
func (rule *Disallow) Validate(string) (bool, error) {
	return false, nil
}

func (rule *Disallow) getMessage() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.Message
}

func (rule *Disallow) GetErrorMessage() string {
	var disallowMessage string
	if rule.getMessage() != "" {
		disallowMessage = fmt.Sprintf(" (%s)", rule.getMessage())
	} else {
		disallowMessage = ""
	}
	return fmt.Sprintf("%s%s", rule.GetName(), disallowMessage)
}
