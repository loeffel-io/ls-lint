package main

import "sync"

type Error struct {
	Path  string
	Rules []*Rule
	*sync.RWMutex
}

func (error *Error) getPath() string {
	error.RLock()
	defer error.RUnlock()

	return error.Path
}

func (error *Error) getRules() []*Rule {
	error.RLock()
	defer error.RUnlock()

	return error.Rules
}

func (error *Error) addRule(rule *Rule) {
	error.Lock()
	defer error.Unlock()

	error.Rules = append(error.Rules, rule)
}
