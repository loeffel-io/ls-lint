package main

import "sync"

type Rule struct {
	Name  string
	Regex string
	*sync.RWMutex
}

func (rule *Rule) getName() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.Name
}

func (rule *Rule) getRegex() string {
	rule.RLock()
	defer rule.RUnlock()

	return rule.Regex
}
