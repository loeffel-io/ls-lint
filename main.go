package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

func main() {
	os.Exit(Run(os.Stdout, os.Args))
}

func Run(writer io.Writer, args []string) int {
	var exitCode = 0
	var flags = flag.NewFlagSet(args[0], flag.ExitOnError)
	flag.CommandLine.SetOutput(writer)
	var warn = flags.Bool("warn", false, "treat lint errors as warnings; write output to stdout and return exit code 0")
	var debug = flags.Bool("debug", false, "write debug informations to stdout")
	var pwd = flags.String("pwd", "", "relative path to the desired working directory")
	var config_file = flags.String("config", "", "relative path to a config file")

	if err := flags.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if *pwd == "" {
		*pwd = cwd
	} else {
		*pwd = path.Join(cwd, *pwd)
	}

	var filesystem = os.DirFS(*pwd)

	var config = &Config{
		RWMutex: new(sync.RWMutex),
	}

	var linter = &Linter{
		Statistic: nil,
		Errors:    make([]*Error, 0),
		RWMutex:   new(sync.RWMutex),
	}

	if err := read_config_file(*pwd, *config_file, config); err != nil {
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
		return exitCode
	}

	if !*warn {
		writer = os.Stderr
		exitCode = 1
	}

	linter.printErrors(writer, cwd, *pwd)

	return exitCode
}
