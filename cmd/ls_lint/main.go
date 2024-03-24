package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	_flag "github.com/loeffel-io/ls-lint/v2/internal/flag"
	"github.com/loeffel-io/ls-lint/v2/internal/linter"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"maps"
	"os"
	"runtime"
	"slices"
	"strings"
	"testing/fstest"
)

var Version = "dev"

func main() {
	var err error
	var exitCode = 0
	var writer = os.Stdout
	var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var flagWorkdir = flags.String("workdir", ".", "change working directory before executing the given subcommand")
	var flagErrorOutputFormat = flags.String("error-output-format", "text", "use a specific error output format (text, json)")
	var flagWarn = flags.Bool("warn", false, "write lint errors to stdout instead of stderr (exit 0)")
	var flagDebug = flags.Bool("debug", false, "write debug informations to stdout")
	var flagVersion = flags.Bool("version", false, "prints version information for ls-lint")

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
		flagConfig = _flag.Config{".ls-lint.yml"}
	}

	var args = flags.Args()
	var filesystem fs.FS
	switch len(args) {
	case 0:
		filesystem = os.DirFS(*flagWorkdir)
	default:
		var mapFilesystem = make(fstest.MapFS, len(args))
		for _, file := range args {
			var fileInfo os.FileInfo
			if fileInfo, err = os.Stat(fmt.Sprintf("%s/%s", *flagWorkdir, file)); err != nil {
				log.Fatal(err)
			}

			mapFilesystem[file] = &fstest.MapFile{Mode: fileInfo.Mode()}
		}
		filesystem = mapFilesystem
	}

	var lslintConfig = config.NewConfig(make(config.Ls), make([]string, 0))
	for _, c := range flagConfig {
		var tmpLslintConfig = config.NewConfig(nil, nil)
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

	var lslintLinter = linter.NewLinter(
		".",
		lslintConfig,
		debug.NewStatistic(),
		make([]*rule.Error, 0),
	)

	if err = lslintLinter.Run(filesystem, *flagDebug); err != nil {
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
		var errIndex = make(map[string][]string, len(lslintLinter.GetErrors()))
		for _, ruleErr := range lslintLinter.GetErrors() {
			errIndex[ruleErr.GetPath()] = make([]string, len(ruleErr.GetRules()))
			for i, ruleErrMessages := range ruleErr.GetRules() {
				errIndex[ruleErr.GetPath()][i] = ruleErrMessages.GetErrorMessage()
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

			for _, errRule := range ruleErr.GetRules() {
				ruleMessages = append(ruleMessages, errRule.GetErrorMessage())
			}

			if _, err = fmt.Fprintf(writer, "%s failed for rules: %s\n", ruleErr.GetPath(), strings.Join(ruleMessages, "|")); err != nil {
				log.Fatal(err)
			}
		}
	}

	os.Exit(exitCode)
}
