package main

import (
	"log"
	"os"
	"sync"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var tests = []*readConfigFileTestCase{
		{name: "Default Case", cwd: cwd, config_file: ""},
		{name: "Specify -config", cwd: cwd, config_file: "./examples/nuxt-nuxt-js/.ls-lint.yml"},
		{name: "Non-existent cwd", cwd: "pathThat-should_never.exist/", config_file: ".ls-lint.yml", expected_error: "open pathThat-should_never.exist/: no such file or directory"},
		{name: "Non-existing config_file", cwd: cwd, config_file: "non-existent-config.yml", expected_error: "open non-existent-config.yml: no such file or directory"},
		{name: "Non-existent config_file without specifying -config", cwd: "./npm", config_file: "", expected_error: "no config file (.ls-lint.yml or .ls-lint.yaml) was found"},
	}
	
	var i = 0
	for _, test := range tests {
		test_config := &Config{
			RWMutex: new(sync.RWMutex),
		}

		err := read_config_file(test.cwd, test.config_file, test_config)

		if err != nil && err.Error() != test.expected_error {
			t.Errorf("Test %d failed (%s) with unmatched error - %s", i, test.name, err.Error())
			return
		}

		i++
	}
}