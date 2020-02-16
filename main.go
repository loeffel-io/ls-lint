package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	var config = &Config{
		RWMutex: new(sync.RWMutex),
	}

	var linter = &Linter{
		Entrypoint: ".",
		Errors:     make([]*Error, 0),
		RWMutex:    new(sync.RWMutex),
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
	err = yaml.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatal(err)
	}

	if err := linter.Run(config); err != nil {
		log.Fatal(err)
	}
}
