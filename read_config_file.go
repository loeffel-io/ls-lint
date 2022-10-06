package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path"
)

func read_config_file(cwd string, config_file string, config *Config) {

	// if config_file is default, see if there's an `.ls-lint.yaml` file instead
	if config_file == ".ls-lint.yml" {
		files, err := os.ReadDir(cwd)
		if err != nil {
			log.Fatal(err)
		}
	
		for _, file := range files {
			match, match_err := path.Match(".ls-lint.y*ml", file.Name())
			if match_err != nil {
				log.Fatal(match_err)
			}
			if match {
				config_file = file.Name()
				break
			}
		}
	}
	// open config file
	file, err := os.Open(config_file)

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
	configBytes, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	// to yaml
	err = yaml.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatal(err)
	}
}