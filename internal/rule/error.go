package rule

import "sync"

type Error struct {
	Path  string
	Rules []Rule
	*sync.RWMutex
}

func (error *Error) GetPath() string {
	error.RLock()
	defer error.RUnlock()

	return error.Path
}

func (error *Error) GetRules() []Rule {
	error.RLock()
	defer error.RUnlock()

	return error.Rules
}
