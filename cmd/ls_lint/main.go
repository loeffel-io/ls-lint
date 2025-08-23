package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	_flag "github.com/loeffel-io/ls-lint/v2/internal/flag"
	"github.com/loeffel-io/ls-lint/v2/internal/linter"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"go.yaml.in/yaml/v3"
)

var Version = "dev"

const (
	// expected ls-lint config file
	lsLintConfigFile = ".ls-lint.yaml"
	// former ls-lint config file, supported for backward compatibility
	lsLintConfigFileLegacy = ".ls-lint.yml"
)

func main() {
	var err error
	exitCode := 0
	writer := os.Stdout
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagWorkdir := flags.String("workdir", ".", "change working directory before executing the given subcommand")
	flagErrorOutputFormat := flags.String("error-output-format", "text", "use a specific error output format (text, json)")
	flagWarn := flags.Bool("warn", false, "write lint errors to stdout instead of stderr (exit 0)")
	flagDebug := flags.Bool("debug", false, "write debug informations to stdout")
	flagVersion := flags.Bool("version", false, "prints version information for ls-lint")

	var flagConfig _flag.Config
	flags.Var(&flagConfig, "config", "ls-lint config file path(s)")

	flags.Usage = func() {
		if _, err = fmt.Fprintln(flags.Output(), "ls-lint [options] [file|dir]*"); err != nil {
			log.Fatal(err)
		}

		if _, err = fmt.Fprintln(flags.Output(), "Options: "); err != nil {
			log.Fatal(err)
		}

		flags.PrintDefaults()
	}

	if err = flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *flagVersion {
		fmt.Printf("ls-lint %s\ngo %s\n", Version, runtime.Version())
		os.Exit(0)
	}

	if len(flagConfig) == 0 {
		// no config files was provided by the --config flag

		// We try the .yaml file first
		configFile := lsLintConfigFile
		if _, err := os.Stat(lsLintConfigFileLegacy); err == nil {
			// but we use the .yml one if it exists
			configFile = lsLintConfigFileLegacy
		}
		flagConfig = _flag.Config{configFile}
	}

	filesystem := os.DirFS(*flagWorkdir)
	var paths map[string]struct{}
	if len(flags.Args()[0:]) > 0 {
		paths = make(map[string]struct{}, len(flags.Args()[0:]))
		for _, path := range flags.Args()[0:] {
			paths[path] = struct{}{}
		}
	}

	lslintConfig := config.NewConfig(make(config.Ls), make([]string, 0))
	for _, c := range flagConfig {
		tmpLslintConfig := config.NewConfig(nil, nil)
		var tmpConfigBytes []byte

		if tmpConfigBytes, err = os.ReadFile(c); err != nil {
			log.Fatal(err)
		}

		if err = yaml.Unmarshal(tmpConfigBytes, tmpLslintConfig); err != nil {
			log.Fatal(err)
		}

		maps.Copy(lslintConfig.GetLs(), tmpLslintConfig.GetLs())
		lslintConfig.Ignore = append(lslintConfig.Ignore, tmpLslintConfig.GetIgnore()...)
		slices.Sort(lslintConfig.Ignore)
		lslintConfig.Ignore = slices.Compact(lslintConfig.Ignore)
	}

	lslintLinter := linter.NewLinter(
		".",
		lslintConfig,
		debug.NewStatistic(),
		make([]*rule.Error, 0),
	)

	if err = lslintLinter.Run(filesystem, paths, *flagDebug); err != nil {
		log.Fatal(err)
	}

	ruleErrors := lslintLinter.GetErrors()

	if len(ruleErrors) == 0 {
		os.Exit(exitCode)
	}

	if !*flagWarn {
		writer = os.Stderr
		exitCode = 1
	}

	switch *flagErrorOutputFormat {
	case "json":
		errIndex := make(map[string]map[string][]string, len(lslintLinter.GetErrors()))
		for _, ruleErr := range lslintLinter.GetErrors() {
			path := ruleErr.GetPath()
			if path == "" {
				path = "."
			}

			if _, ok := errIndex[path]; !ok {
				errIndex[path] = make(map[string][]string)
			}

			for _, errRule := range ruleErr.GetRules() {
				if !ruleErr.IsDir() && errRule.GetName() == "exists" {
					continue
				}

				errIndex[path][ruleErr.GetExt()] = append(errIndex[path][ruleErr.GetExt()], errRule.GetErrorMessage())
			}
		}

		var jsonStr []byte
		if jsonStr, err = json.Marshal(errIndex); err != nil {
			log.Fatal(err)
		}

		if _, err = fmt.Fprintln(writer, string(jsonStr)); err != nil {
			log.Fatal(err)
		}
	default:
		for _, ruleErr := range lslintLinter.GetErrors() {
			var ruleMessages []string

			path := ruleErr.GetPath()
			if path == "" {
				path = "."
			}

			for _, errRule := range ruleErr.GetRules() {
				if !ruleErr.IsDir() && errRule.GetName() == "exists" {
					continue
				}

				ruleMessages = append(ruleMessages, errRule.GetErrorMessage())
			}

			if _, err = fmt.Fprintf(writer, "%s failed for `%s` rules: %s\n", path, ruleErr.GetExt(), strings.Join(ruleMessages, " | ")); err != nil {
				log.Fatal(err)
			}
		}
	}

	os.Exit(exitCode)
}
