package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path"
)

func read_config_file(cwd string, config_file string, config *Config) error {

	// if config_file is empty, check for both `.ls-lint.yml` and `.ls-lint.yaml`
	if config_file == "" {
		files, err := os.ReadDir(cwd)
		if err != nil {
			return err
		}
	
		for _, file := range files {
			match, match_err := path.Match(".ls-lint.y*ml", file.Name())
			if match_err != nil {
				return match_err
			}
			if match {
				config_file = file.Name()
				break
			}
		}
	}
	if config_file == "" {
		return errors.New("no config file (.ls-lint.yml or .ls-lint.yaml) was found")
	}
	// open config file
	file, err := os.Open(config_file)

	if err != nil {
		return err
	}

	// close file
	defer func() {
		err = file.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	// read file
	configBytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	// to yaml
	err = yaml.Unmarshal(configBytes, &config)

	return err
}
