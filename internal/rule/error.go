package rule

import "sync"

type Error struct {
	Path  string
	Dir   bool
	Ext   string
	Rules []Rule
	*sync.RWMutex
}

func (err *Error) GetPath() string {
	err.RLock()
	defer err.RUnlock()

	return err.Path
}

func (err *Error) IsDir() bool {
	err.RLock()
	defer err.RUnlock()

	return err.Dir
}

func (err *Error) GetExt() string {
	err.RLock()
	defer err.RUnlock()

	return err.Ext
}

func (err *Error) GetRules() []Rule {
	err.RLock()
	defer err.RUnlock()

	return err.Rules
}
