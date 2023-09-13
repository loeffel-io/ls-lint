package main

import (
	"flag"
	"fmt"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	_flag "github.com/loeffel-io/ls-lint/v2/internal/flag"
	"github.com/loeffel-io/ls-lint/v2/internal/linter"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"gopkg.in/yaml.v3"
	"log"
	"maps"
	"os"
	"runtime"
	"slices"
	"strings"
)

var Version = "dev"

func main() {
	var err error
	var exitCode = 0
	var writer = os.Stdout
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var flagWorkdir = flags.String("workdir", ".", "change working directory before executing the given subcommand")
	var flagWarn = flags.Bool("warn", false, "write lint errors to stdout instead of stderr (exit 0)")
	var flagDebug = flags.Bool("debug", false, "write debug informations to stdout")
	var flagVersion = flags.Bool("version", false, "prints version information for ls-lint")

	var flagConfig _flag.Config
	flags.Var(&flagConfig, "config", "ls-lint config file path(s)")

	if err = flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *flagVersion {
		fmt.Printf("ls-lint %s\ngo %s\n", Version, runtime.Version())
		os.Exit(0)
	}

	var filesystem = os.DirFS(*flagWorkdir)

	if len(flagConfig) == 0 {
		flagConfig = _flag.Config{".ls-lint.yaml"}
	}

	var lslintConfig = config.NewConfig(make(config.Ls), make([]string, 0))
	for _, c := range flagConfig {
		var tmpLslintConfig = config.NewConfig(nil, nil)
		var tmpConfigBytes []byte

		// read file
		if tmpConfigBytes, err = os.ReadFile(c); err != nil {
			log.Fatal(err)
		}

		// to yaml
		if err = yaml.Unmarshal(tmpConfigBytes, tmpLslintConfig); err != nil {
			log.Fatal(err)
		}

		maps.Copy(lslintConfig.GetLs(), tmpLslintConfig.GetLs())
		lslintConfig.Ignore = append(lslintConfig.Ignore, tmpLslintConfig.GetIgnore()...)
		slices.Sort(lslintConfig.Ignore)
		lslintConfig.Ignore = slices.Compact(lslintConfig.Ignore)
	}

	// linter
	var lslintLinter = linter.NewLinter(
		".",
		lslintConfig,
		debug.NewStatistic(),
		make([]*rule.Error, 0),
	)

	// runner
	if err = lslintLinter.Run(filesystem, *flagDebug); err != nil {
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
