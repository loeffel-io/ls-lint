package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func getFullPath(path string) string {
	return fmt.Sprintf("%s%s%s", root, sep, path)
}

func main() {
	var filesystem = os.DirFS(root)

	var config = &Config{
		RWMutex: new(sync.RWMutex),
	}

	var linter = &Linter{
		Errors:  make([]*Error, 0),
		RWMutex: new(sync.RWMutex),
	}

	// open config file
	file, err := os.Open(".ls-lint.yml")

	if err != nil {
		log.Fatal(err)
	}

	// close file
	defer func() {
		err = file.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	// read file
	configBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	// to yaml
	err = yaml.Unmarshal(normalizeConfig(configBytes, byte(runeUnixSep), byte(runeSep)), &config)

	if err != nil {
		log.Fatal(err)
	}

	// runner
	if err := linter.Run(filesystem, config); err != nil {
		log.Fatal(err)
	}

	// errors
	errors := linter.getErrors()

	// no errors
	if len(errors) == 0 {
		os.Exit(0)
	}

	// with errors
	for _, err := range linter.getErrors() {
		var ruleMessages []string

		for _, rule := range err.getRules() {
			ruleMessages = append(ruleMessages, rule.GetErrorMessage())
		}

		log.Printf("%s failed for rules: %s", err.getPath(), strings.Join(ruleMessages, fmt.Sprintf(" %s ", or)))
	}

	os.Exit(1)
}
