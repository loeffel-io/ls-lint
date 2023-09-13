package rule

import (
	"sync"
	"unicode"
)

type ScreamingSnakeCase struct {
	Name string
	*sync.RWMutex
}

func (rule *ScreamingSnakeCase) Init() Rule {
	rule.Name = "screamingsnakecase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *ScreamingSnakeCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *ScreamingSnakeCase) SetParameters(params []string) error {
	return nil
}

func (rule *ScreamingSnakeCase) GetParameters() []string {
	return nil
}

// Validate checks if string is screaming sneak case
// false if rune is no uppercase letter, digit or _
func (rule *ScreamingSnakeCase) Validate(value string) (bool, error) {
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
