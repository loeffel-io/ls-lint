package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	var exitCode = 0
	var writer = os.Stdout
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var warn = flags.Bool("warn", false, "treat lint errors as warnings; write output to stdout and return exit code 0")
	var debug = flags.Bool("debug", false, "write debug informations to stdout")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	var filesystem = os.DirFS(root)

	var config = &Config{
		RWMutex: new(sync.RWMutex),
	}

	var linter = &Linter{
		Statistic: nil,
		Errors:    make([]*Error, 0),
		RWMutex:   new(sync.RWMutex),
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
	if err := linter.Run(filesystem, config, *debug, false); err != nil {
		log.Fatal(err)
	}

	// errors
	errors := linter.getErrors()

	// no errors
	if len(errors) == 0 {
		os.Exit(exitCode)
	}

	if !*warn {
		writer = os.Stderr
		exitCode = 1
	}

	logger := log.New(writer, "", log.LstdFlags)

	// with errors
	for _, err := range linter.getErrors() {
		var ruleMessages []string

		for _, rule := range err.getRules() {
			ruleMessages = append(ruleMessages, rule.GetErrorMessage())
		}

		logger.Printf("%s failed for rules: %s", err.getPath(), strings.Join(ruleMessages, fmt.Sprintf(" %s ", or)))
	}

	os.Exit(exitCode)
}
