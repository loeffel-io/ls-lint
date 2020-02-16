package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Linter struct {
	Entrypoint string
	Errors     []*Error
	*sync.RWMutex
}

func (linter *Linter) getEntrypoint() string {
	linter.RLock()
	defer linter.RUnlock()

	return linter.Entrypoint
}

func (linter *Linter) getErrors() []*Error {
	linter.RLock()
	defer linter.RUnlock()

	return linter.Errors
}

func (linter *Linter) addError(error *Error) {
	linter.Lock()
	defer linter.Unlock()

	linter.Errors = append(linter.Errors, error)
}

func (linter *Linter) Run(config *Config) error {
	var g = new(errgroup.Group)
	var index = config.getIndex()

	g.Go(func() error {
		return filepath.Walk(linter.getEntrypoint(), func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				log.Printf("%+v", config.getConfig(index, path))
			}

			return nil
		})
	})

	return g.Wait()
}
