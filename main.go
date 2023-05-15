package main

import (
	"flag"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/linter"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	var exitCode = 0
	var writer = os.Stdout
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var flagConfig = flags.String("config", ".ls-lint.yml", "ls-lint config file path")
	var flagChdir = flags.String("chdir", ".", "Switch to a different working directory before executing the given subcommand")
	var flagWarn = flags.Bool("warn", false, "treat lint ruleErrors as warnings; write output to stdout and return exit code 0")
	var flagDebug = flags.Bool("debug", false, "write debug informations to stdout")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	var filesystem = os.DirFS(*flagChdir)

	var lslintConfig = &config.Config{
		RWMutex: new(sync.RWMutex),
	}

	// read file
	configBytes, err := os.ReadFile(*flagConfig)

	if err != nil {
		log.Fatal(err)
	}

	// to yaml
	err = yaml.Unmarshal(configBytes, &lslintConfig)

	if err != nil {
		log.Fatal(err)
	}

	// linter
	var lslintLinter = linter.NewLinter(
		*flagChdir,
		lslintConfig,
		nil,
		make([]*rule.Error, 0),
	)

	// runner
	if err = lslintLinter.Run(filesystem, *flagDebug, false); err != nil {
		log.Fatal(err)
	}

	// rule errors
	ruleErrors := lslintLinter.GetErrors()

	// no ruleErrors
	if len(ruleErrors) == 0 {
		os.Exit(exitCode)
	}

	if !*flagWarn {
		writer = os.Stderr
		exitCode = 1
	}

	logger := log.New(writer, "", log.LstdFlags)

	// with rule errors
	for _, ruleErr := range lslintLinter.GetErrors() {
		var ruleMessages []string

		for _, errRule := range ruleErr.GetRules() {
			ruleMessages = append(ruleMessages, errRule.GetErrorMessage())
		}

		logger.Printf("%s failed for rules: %s", ruleErr.GetPath(), strings.Join(ruleMessages, "|"))
	}

	os.Exit(exitCode)
}
