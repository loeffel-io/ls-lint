package main

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"os"
	"reflect"
	"strings"
	"sync"
)

type ls map[interface{}]interface{}
type index map[string]map[string][]Rule

const (
	sep    = string(os.PathSeparator)
	extSep = "."
	root   = "."
	dir    = ".dir"
	or     = "|"
)

type Config struct {
	Ls     ls       `yaml:"ls"`
	Ignore []string `yaml:"ignore"`
	*sync.RWMutex
}

func (config *Config) getLs() ls {
	config.RLock()
	defer config.RUnlock()

	return config.shiftLs(config.Ls)
}

func (config *Config) shiftLs(list ls) ls {
	var shift = make(ls)
	shift[root] = list

	return shift
}

func (config *Config) getIgnore() []string {
	config.RLock()
	defer config.RUnlock()

	return config.Ignore
}

func (config *Config) getIgnoreIndex() map[string]bool {
	var ignoreIndex = make(map[string]bool)

	for _, path := range config.getIgnore() {
		ignoreIndex[path] = true
		ignoreIndex[fmt.Sprintf("%s%s%s", root, sep, path)] = true
	}

	return ignoreIndex
}

func (config *Config) shouldIgnore(ignoreIndex map[string]bool, path string) bool {
	if ignore, exists := ignoreIndex[path]; exists {
		return ignore
	}

	if ignore, exists := ignoreIndex[getFullPath(path)]; exists {
		return ignore
	}

	dirs := strings.Split(path, sep)
	for i := 0; i < len(dirs); i++ {
		if ignore, exists := ignoreIndex[strings.Join(dirs[:i], sep)]; exists {
			return ignore
		}

		if ignore, exists := ignoreIndex[getFullPath(strings.Join(dirs[:i], sep))]; exists {
			return ignore
		}
	}

	return false
}

func (config *Config) getConfig(index index, path string) map[string][]Rule {
	dirs := strings.Split(path, sep)

	for i := len(dirs); i >= 0; i-- {
		if find, exists := index[strings.Join(dirs[:i], sep)]; exists {
			return find
		}
	}

	return nil
}

func (config *Config) copyRule(rule Rule) Rule {
	switch rule.GetName() {
	case "regex":
		return new(RuleRegex).Init()
	}

	return rule
}

func (config *Config) walkIndex(index index, key string, list ls) error {
	if index[key] == nil {
		index[key] = make(map[string][]Rule)
	}

	for k, v := range list {
		if v == nil {
			continue
		}

		if reflect.TypeOf(v).Kind() == reflect.Map {
			if err := config.walkIndex(index, fmt.Sprintf("%s%s%s", key, sep, k.(string)), v.(ls)); err != nil {
				return err
			}

			continue
		}

		for _, ruleName := range strings.Split(v.(string), or) {
			ruleName = strings.TrimSpace(ruleName)
			ruleSplit := strings.SplitN(ruleName, ":", 2)
			ruleName = ruleSplit[0]

			if rule, exists := rules[ruleName]; exists {
				rule = config.copyRule(rule)

				if err := rule.SetParameters(ruleSplit[1:]); err != nil {
					return fmt.Errorf("rule %s failed with %s", ruleName, err.Error())
				}

				index[key][k.(string)] = append(index[key][k.(string)], rule)
				continue
			}

			return fmt.Errorf("rule %s not exists", ruleName)
		}
	}

	return nil
}

func (config *Config) getIndex(list ls) (index, error) {
	var index = make(index)

	for key, value := range list {
		if err := config.walkIndex(index, key.(string), value.(ls)); err != nil {
			return nil, err
		}
	}

	return index, nil
}

func (config *Config) globIndex(index index) (err error) {
	for key, value := range index {
		var matches []string

		if !strings.ContainsAny(key, "*{}") {
			continue
		}

		if matches, err = doublestar.Glob(key); err != nil {
			return err
		}

		if len(matches) == 0 {
			delete(index, key)
			continue
		}

		for _, match := range matches {
			var matchInfo os.FileInfo
			match = getFullPath(match)

			if matchInfo, err = os.Stat(match); err != nil {
				return err
			}

			if !matchInfo.IsDir() {
				continue
			}

			index[match] = value
			delete(index, key)
		}
	}

	return nil
}
