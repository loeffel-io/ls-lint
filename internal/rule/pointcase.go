package rule

import (
	"sync"
	"unicode"
)

type PointCase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *PointCase) Init() Rule {
	rule.name = "pointcase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *PointCase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *PointCase) SetParameters(params []string) error {
	return nil
}

func (rule *PointCase) GetParameters() []string {
	return nil
}

func (rule *PointCase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if string is "point case"
// false if rune is no lowercase letter, digit or .
func (rule *PointCase) Validate(value string, fail bool) (bool, error) {
	for _, c := range value {
		if c == 46 || unicode.IsDigit(c) { // 46 => .
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

func (rule *PointCase) GetErrorMessage() string {
	return rule.GetName()
}
