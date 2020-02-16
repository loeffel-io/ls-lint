package main

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type ls map[string]interface{}
type index map[string]map[string]string

type Config struct {
	Ls ls `yaml:"ls"`
	*sync.RWMutex
}

func (config *Config) getLs() ls {
	config.RLock()
	defer config.RUnlock()

	return config.Ls
}

func (config *Config) getConfig(index index, path string) map[string]string {
	const sep = "/"
	dirs := strings.Split(path, sep)

	for i := len(dirs); i >= 0; i-- {
		if find, exists := index[strings.Join(dirs[:i], sep)]; exists {
			return find
		}
	}

	return nil
}

func (config *Config) walkIndex(index index, key string, value map[interface{}]interface{}) {
	const sep = "/"

	if index[key] == nil {
		index[key] = make(map[string]string)
	}

	for k, v := range value {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			config.walkIndex(index, fmt.Sprintf("%s%s%s", key, sep, k.(string)), v.(map[interface{}]interface{}))
			continue
		}

		index[key][k.(string)] = v.(string)
	}
}

func (config *Config) getIndex(ls ls) index {
	var index = make(index)

	for key, value := range ls {
		config.walkIndex(index, key, value.(map[interface{}]interface{}))
	}

	return index
}
