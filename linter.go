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
	var (
		g     = new(errgroup.Group)
		ls    = config.getLs()
		index = config.getIndex(ls)
	)

	for entrypoint := range ls {
		g.Go(func() error {
			return filepath.Walk(entrypoint, func(path string, info os.FileInfo, err error) error {
				log.Printf("%+v", config.getConfig(index, path))

				if info.IsDir() {
				}

				return nil
			})
		})
	}

	return g.Wait()
}
