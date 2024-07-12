package rule

import (
	"sync"
	"unicode"
)

type SnakeCase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *SnakeCase) Init() Rule {
	rule.name = "snakecase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *SnakeCase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *SnakeCase) SetParameters(params []string) error {
	return nil
}

func (rule *SnakeCase) GetParameters() []string {
	return nil
}

func (rule *SnakeCase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if string is sneak case
// false if rune is no lowercase letter, digit or _
func (rule *SnakeCase) Validate(value string, fail bool) (bool, error) {
	for _, c := range value {
		if c == 95 || unicode.IsDigit(c) { // 95 => _
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

func (rule *SnakeCase) GetErrorMessage() string {
	return rule.GetName()
}

func (rule *SnakeCase) Copy() Rule {
	return rule
}
