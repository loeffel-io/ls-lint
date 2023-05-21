package main

import (
	"flag"
	"fmt"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	"github.com/loeffel-io/ls-lint/v2/internal/linter"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"runtime"
	"strings"
)

var Version = "dev"

func main() {
	var err error
	var exitCode = 0
	var writer = os.Stdout
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var flagConfig = flags.String("config", ".ls-lint.yml", "ls-lint config file path")
	var flagWorkdir = flags.String("workdir", ".", "change working directory before executing the given subcommand")
	var flagWarn = flags.Bool("warn", false, "write lint errors to stdout instead of stderr (exit 0)")
	var flagDebug = flags.Bool("debug", false, "write debug informations to stdout")
	var flagVersion = flags.Bool("version", false, "prints version information for ls-lint")

	if err = flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *flagVersion {
		fmt.Printf("ls-lint %s\ngo %s\n", Version, runtime.Version())
		os.Exit(0)
	}

	var filesystem = os.DirFS(*flagWorkdir)
	var lslintConfig = config.NewConfig(nil, nil)
	var configBytes []byte

	// read file
	if configBytes, err = os.ReadFile(*flagConfig); err != nil {
		log.Fatal(err)
	}

	// to yaml
	if err = yaml.Unmarshal(configBytes, lslintConfig); err != nil {
		log.Fatal(err)
	}

	// linter
	var lslintLinter = linter.NewLinter(
		".",
		lslintConfig,
		debug.NewStatistic(),
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
