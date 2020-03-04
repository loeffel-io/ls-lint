package main

import "sync"

type Error struct {
	Path  string
	Rules []Rule
	*sync.RWMutex
}

func (error *Error) getPath() string {
	error.RLock()
	defer error.RUnlock()

	return error.Path
}

func (error *Error) getRules() []Rule {
	error.RLock()
	defer error.RUnlock()

	return error.Rules
}
