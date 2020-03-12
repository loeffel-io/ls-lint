package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Linter struct {
	Errors []*Error
	*sync.RWMutex
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

func (linter *Linter) validateDir(config *Config, index index, path string) error {
	var g = new(errgroup.Group)
	var errRules = make([]Rule, 0)

	rules := config.getConfig(index, path)
	basename := filepath.Base(path)

	if basename == root {
		return nil
	}

	if _, exists := rules[dir]; !exists {
		return nil
	}

	for _, rule := range rules[dir] {
		g.Go(func() error {
			valid, err := rule.Validate(basename)

			if err != nil {
				return err
			}

			if !valid {
				errRules = append(errRules, rule)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if len(errRules) > 0 {
		linter.addError(&Error{
			Path:    path,
			Rules:   errRules,
			RWMutex: new(sync.RWMutex),
		})
	}

	return nil
}

func (linter *Linter) validateFile(config *Config, index index, entrypoint string, path string) error {
	var g = new(errgroup.Group)
	var errRules = make([]Rule, 0)

	exts := strings.Split(filepath.Base(path), extSep)
	rules := config.getConfig(index, path)

	for i := 1; i < len(exts); i++ {
		ext := fmt.Sprintf("%s%s", extSep, strings.Join(exts[i:], extSep))
		withoutExt := strings.TrimSuffix(filepath.Base(path), ext)

		if _, exists := rules[ext]; exists {
			for _, rule := range rules[ext] {
				g.Go(func() error {
					valid, err := rule.Validate(withoutExt)

					if err != nil {
						return err
					}

					if !valid {
						errRules = append(errRules, rule)
					}

					return nil
				})
			}

			break
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if len(errRules) > 0 {
		linter.addError(&Error{
			Path:    path,
			Rules:   errRules,
			RWMutex: new(sync.RWMutex),
		})
	}

	return nil
}

func (linter *Linter) Run(config *Config) error {
	var g = new(errgroup.Group)
	var ls = config.getLs()
	var ignoreIndex = config.getIgnoreIndex()
	var index, err = config.getIndex(ls)

	if err != nil {
		return err
	}

	for entrypoint := range ls {
		g.Go(func() error {
			return filepath.Walk(entrypoint.(string), func(path string, info os.FileInfo, err error) error {
				if config.shouldIgnore(ignoreIndex, path) {
					return nil
				}

				path = getFullPath(path)

				if info == nil {
					return fmt.Errorf("%s not found", entrypoint)
				}

				if info.IsDir() {
					return linter.validateDir(config, index, path)
				}

				return linter.validateFile(config, index, entrypoint.(string), path)
			})
		})
	}

	return g.Wait()
}
