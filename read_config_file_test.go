package main

import (
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	examples_config := &Config{
		RWMutex: new(sync.RWMutex),
	}
	root_config := &Config {
		RWMutex: new(sync.RWMutex),
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	read_config_file(cwd, "./examples/nuxt-nuxt-js/.ls-lint.yml", examples_config)
	read_config_file(cwd, ".ls-lint.yml", root_config)
	if reflect.DeepEqual(examples_config.getLs(), root_config.getLs()) {
		t.Errorf("Both configuration files have the same ls rules")
	}
}