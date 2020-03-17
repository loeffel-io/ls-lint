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
	var rulesError = 0
	var rulesErrorMutex = new(sync.Mutex)

	rules := config.getConfig(index, path)
	basename := filepath.Base(path)

	if basename == root {
		return nil
	}

	if _, exists := rules[dir]; !exists {
		return nil
	}

	for _, rule := range rules[dir] {
		rule := rule
		g.Go(func() error {
			valid, err := rule.Validate(basename)

			if err != nil {
				return err
			}

			if !valid {
				rulesErrorMutex.Lock()
				rulesError += 1
				rulesErrorMutex.Unlock()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if rulesError == 0 || rulesError != len(rules[dir]) {
		return nil
	}

	linter.addError(&Error{
		Path:    path,
		Rules:   rules[dir],
		RWMutex: new(sync.RWMutex),
	})

	return nil
}

func (linter *Linter) validateFile(config *Config, index index, entrypoint string, path string) error {
	var ext string
	var g = new(errgroup.Group)
	var rulesError = 0
	var rulesErrorMutex = new(sync.Mutex)

	exts := strings.Split(filepath.Base(path), extSep)
	rules := config.getConfig(index, path)

	for i := 1; i < len(exts); i++ {
		ext = fmt.Sprintf("%s%s", extSep, strings.Join(exts[i:], extSep))
		withoutExt := strings.TrimSuffix(filepath.Base(path), ext)

		if _, exists := rules[ext]; exists {
			for _, rule := range rules[ext] {
				rule := rule
				g.Go(func() error {
					valid, err := rule.Validate(withoutExt)

					if err != nil {
						return err
					}

					if !valid {
						rulesErrorMutex.Lock()
						rulesError += 1
						rulesErrorMutex.Unlock()
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

	if rulesError == 0 || rulesError != len(rules[ext]) {
		return nil
	}

	linter.addError(&Error{
		Path:    path,
		Rules:   rules[ext],
		RWMutex: new(sync.RWMutex),
	})

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
		entrypoint := entrypoint
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
