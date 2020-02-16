package main

import (
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
	dirs := strings.Split(path, "/")

	for i := len(dirs); i >= 0; i-- {
		if find, exists := index[strings.Join(dirs[:i], "/")]; exists {
			return find
		}
	}

	return nil
}

func (config *Config) getIndex() index {
	return index{
		"src/js": {
			"test": "asdf",
		},
	}
}
