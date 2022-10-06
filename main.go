package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

func main() {
	var exitCode = 0
	var writer = os.Stdout
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var warn = flags.Bool("warn", false, "treat lint errors as warnings; write output to stdout and return exit code 0")
	var debug = flags.Bool("debug", false, "write debug informations to stdout")
	var config_file = flags.String("config", "", "relative path to a config file, its directory is the new root")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var filesystem = os.DirFS(path.Join(cwd, path.Dir(*config_file)))

	var config = &Config{
		RWMutex: new(sync.RWMutex),
	}

	var linter = &Linter{
		Statistic: nil,
		Errors:    make([]*Error, 0),
		RWMutex:   new(sync.RWMutex),
	}
	
	if err := read_config_file(cwd, *config_file, config); err != nil {
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
