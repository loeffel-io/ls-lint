package rule

import (
	"fmt"
	"sync"
)

type Required struct {
	name      string
	exclusive bool
	value     string
	count     uint16
	*sync.RWMutex
}

func (rule *Required) Init() Rule {
	rule.name = "required"
	rule.exclusive = true
	rule.count = 0
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Required) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *Required) SetParameters(params []string) error {
	rule.Lock()
	defer rule.Unlock()

	if len(params) == 0 {
		return fmt.Errorf("required value not exists")
	}

	if params[0] == "" {
		return fmt.Errorf("required value is empty")
	}

	rule.value = params[0]

	return nil
}

func (rule *Required) GetParameters() []string {
	return []string{rule.getValue()}
}

func (rule *Required) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

func (rule *Required) Validate(value string, _ string, fail bool) (bool, error) {
	if !fail {
		if value == rule.getValue() {
			rule.incrementCount()
		}

		return true, nil
	}

	return rule.getCount() > 0, nil
}

func (rule *Required) getValue() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.value
}

func (rule *Required) getCount() uint16 {
	rule.RLock()
	defer rule.RUnlock()

	return rule.count
}

func (rule *Required) incrementCount() {
	rule.Lock()
	defer rule.Unlock()

	rule.count++
}

func (rule *Required) GetErrorMessage() string {
	return fmt.Sprintf("%s:%s (found %d)", rule.GetName(), rule.getValue(), rule.getCount())
}

func (rule *Required) Copy() Rule {
	rule.RLock()
	defer rule.RUnlock()

	c := new(Required)
	c.Init()
	c.value = rule.value

	return c
}
