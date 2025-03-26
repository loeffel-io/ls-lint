package rule

import (
	"sync"
	"unicode"
)

type KebabCase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *KebabCase) Init() Rule {
	rule.name = "kebabcase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *KebabCase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *KebabCase) SetParameters(params []string) error {
	return nil
}

func (rule *KebabCase) GetParameters() []string {
	return nil
}

func (rule *KebabCase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if string is kebab case
// false if rune is no lowercase letter, digit or -
func (rule *KebabCase) Validate(value string, _ string, _ bool) (bool, error) {
	for _, c := range value {
		if c == 45 || unicode.IsDigit(c) { // 45 => -
			continue
		}

		if !unicode.IsLetter(c) {
			return false, nil
		}

		if !unicode.IsLower(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *KebabCase) GetErrorMessage() string {
	return rule.GetName()
}

func (rule *KebabCase) Copy() Rule {
	return rule
}
