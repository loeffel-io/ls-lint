package main

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type ls map[string]interface{}
type index map[string]map[string][]Rule

const sep = "/"

type Config struct {
	Ls ls `yaml:"ls"`
	*sync.RWMutex
}

func (config *Config) getLs() ls {
	config.RLock()
	defer config.RUnlock()

	return config.Ls
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

func (config *Config) walkIndex(index index, key string, value map[interface{}]interface{}) error {
	if index[key] == nil {
		index[key] = make(map[string][]Rule)
	}

	for k, v := range value {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			if err := config.walkIndex(index, fmt.Sprintf("%s%s%s", key, sep, k.(string)), v.(map[interface{}]interface{})); err != nil {
				return err
			}

			continue
		}

		for _, ruleName := range strings.Split(v.(string), ",") {
			ruleName = strings.TrimSpace(ruleName)

			if rule, exists := rules[ruleName]; exists {
				index[key][k.(string)] = append(index[key][k.(string)], rule)
				continue
			}

			return fmt.Errorf("rule %s not exists", ruleName)
		}
	}

	return nil
}

func (config *Config) getIndex(ls ls) (index, error) {
	var index = make(index)

	for key, value := range ls {
		if err := config.walkIndex(index, key, value.(map[interface{}]interface{})); err != nil {
			return nil, err
		}
	}

	return index, nil
}
