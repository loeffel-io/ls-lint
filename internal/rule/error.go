package rule

import "sync"

type Error struct {
	Path  string
	Dir   bool
	Ext   string
	Rules []Rule
	*sync.RWMutex
}

func (error *Error) GetPath() string {
	error.RLock()
	defer error.RUnlock()

	return error.Path
}

func (error *Error) IsDir() bool {
	error.RLock()
	defer error.RUnlock()

	return error.Dir
}

func (error *Error) GetExt() string {
	error.RLock()
	defer error.RUnlock()

	return error.Ext
}

func (error *Error) GetRules() []Rule {
	error.RLock()
	defer error.RUnlock()

	return error.Rules
}
