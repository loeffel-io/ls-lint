package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
	"strings"
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
				if info == nil {
					return fmt.Errorf("%s not found", entrypoint)
				}

				if info.IsDir() {
					rules := index[path]
					basename := filepath.Base(path)

					log.Printf("%+v %s", rules, basename)
					return nil
				}

				ext := filepath.Ext(path)
				rules := index[entrypoint][ext]
				withoutExt := strings.TrimSuffix(filepath.Base(path), ext)

				log.Printf("%s %s %+v", ext, withoutExt, rules)
				return nil
			})
		})
	}

	return g.Wait()
}
