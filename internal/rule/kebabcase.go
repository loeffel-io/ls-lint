package rule

import (
	"sync"
	"unicode"
)

type KebabCase struct {
	Name string
	*sync.RWMutex
}

func (rule *KebabCase) Init() Rule {
	rule.Name = "kebabcase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *KebabCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *KebabCase) SetParameters(params []string) error {
	return nil
}

func (rule *KebabCase) GetParameters() []string {
	return nil
}

// Validate checks if string is kebab case
// false if rune is no lowercase letter, digit or -
func (rule *KebabCase) Validate(value string) (bool, error) {
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
