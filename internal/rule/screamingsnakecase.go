package rule

import (
	"sync"
	"unicode"
)

type ScreamingSnakeCase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *ScreamingSnakeCase) Init() Rule {
	rule.name = "screamingsnakecase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *ScreamingSnakeCase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *ScreamingSnakeCase) SetParameters(params []string) error {
	return nil
}

func (rule *ScreamingSnakeCase) GetParameters() []string {
	return nil
}

func (rule *ScreamingSnakeCase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

// Validate checks if string is screaming sneak case
// false if rune is no uppercase letter, digit or _
func (rule *ScreamingSnakeCase) Validate(value string, fail bool) (bool, error) {
	for _, c := range value {
		if c == 95 || unicode.IsDigit(c) { // 95 => _
			continue
		}

		if !unicode.IsLetter(c) {
			return false, nil
		}

		if !unicode.IsUpper(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *ScreamingSnakeCase) GetErrorMessage() string {
	return rule.GetName()
}
