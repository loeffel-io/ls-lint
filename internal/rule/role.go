package rule

import (
	"fmt"
	"sync"
)

type Role struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *Role) Init() Rule {
	rule.name = "role"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Role) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *Role) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	if len(params) == 0 {
		return fmt.Errorf("role name not exists")
	}

	if params[0] == "" {
		return fmt.Errorf("role name is empty")
	}

	return nil
}

func (rule *Role) GetParameters() []string {
	return nil
}

func (rule *Role) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate is noop
func (rule *Role) Validate(_ string, _ string, _ bool) (bool, error) {
	return true, nil
}

func (rule *Role) GetErrorMessage() string {
	return rule.GetName()
}

func (rule *Role) Copy() Rule {
	return rule
}
