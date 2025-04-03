package rule

import (
	"sync"
	"unicode"
)

// Flatcase lints file names to be made up of only lowercase letters and numbers.
//
// It is a variant of the snake case rule, but it does not allow underscores.
type Flatcase struct {
	name      string
	exclusive bool
	*sync.RWMutex
}

func (rule *Flatcase) Init() Rule {
	rule.name = "flatcase"
	rule.exclusive = false
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *Flatcase) GetName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.name
}

func (rule *Flatcase) SetParameters(params []string) error {
	return nil
}

func (rule *Flatcase) GetParameters() []string {
	return nil
}

func (rule *Flatcase) GetExclusive() bool {
	rule.RLock()
	defer rule.RUnlock()

	return rule.exclusive
}

func (rule *Flatcase) Validate(value string, path string, fail bool) (bool, error) {
	if len(value) == 0 {
		// empty strings are not allowed
		return false, nil
	}

	// let's loop on each rune of the string
	for _, c := range value {
		if unicode.IsDigit(c) {
			// numbers are allowed
			continue
		}

		if !unicode.IsLetter(c) {
			// anything that is not a letter is not allowed
			return false, nil
		}

		if !unicode.IsLower(c) {
			// only lowercase letters are allowed
			return false, nil
		}
	}

	// if we reach this point, the string is valid
	// all characters passed the checks
	return true, nil
}

func (rule *Flatcase) GetErrorMessage() string {
	return rule.GetName()
}

func (rule *Flatcase) Copy() Rule {
	return rule
}
