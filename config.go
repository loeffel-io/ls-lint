package main

import (
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"reflect"
	"strings"
	"sync"
)

type ls map[string]interface{}
type index map[string]map[string][]Rule

const (
	sep    = string('/')
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

	return config.Ls
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
	}

	return ignoreIndex
}

func (config *Config) shouldIgnore(ignoreIndex map[string]bool, path string) bool {
	if ignore, exists := ignoreIndex[path]; exists {
		return ignore
	}

	dirs := strings.Split(path, sep)
	for i := 0; i < len(dirs); i++ {
		if ignore, exists := ignoreIndex[strings.Join(dirs[:i], sep)]; exists {
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
			switch key == "" {
			case true:
				if err := config.walkIndex(index, k, v.(ls)); err != nil {
					return err
				}
			case false:
				var keyCombination = fmt.Sprintf("%s%s%s", key, sep, k)
				if err := config.walkIndex(index, keyCombination, v.(ls)); err != nil {
					return err
				}
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

				index[key][k] = append(index[key][k], rule)
				continue
			}

			return fmt.Errorf("rule %s not exists", ruleName)
		}
	}

	return nil
}

func (config *Config) getIndex(list ls) (index, error) {
	var index = make(index)

	if err := config.walkIndex(index, "", list); err != nil {
		return nil, err
	}

	return index, nil
}

func globIndex[V bool | map[string][]Rule](filesystem fs.FS, index map[string]V) (err error) {
	for key, value := range index {
		var matches []string

		if !strings.ContainsAny(key, "*{}") {
			continue
		}

		if matches, err = doublestar.Glob(filesystem, key); err != nil {
			return err
		}

		if len(matches) == 0 {
			delete(index, key)
			continue
		}

		for _, match := range matches {
			var matchInfo fs.FileInfo

			if matchInfo, err = fs.Stat(filesystem, match); err != nil {
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
