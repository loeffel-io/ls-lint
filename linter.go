package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type Linter struct {
	Statistic *Statistic
	Errors    []*Error
	*sync.RWMutex
}

func (linter *Linter) getStatistic() *Statistic {
	linter.RLock()
	defer linter.RUnlock()

	return linter.Statistic
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

func (linter *Linter) validateFile(config *Config, index index, path string) error {
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

func (linter *Linter) Run(filesystem fs.FS, config *Config, debug bool, statistics bool) (err error) {
	var index index
	var ls = config.getLs()
	var ignoreIndex = config.getIgnoreIndex()

	// create index
	if index, err = config.getIndex(ls); err != nil {
		return err
	}

	// apply globbing to the index
	if err := globIndex(filesystem, index); err != nil {
		return err
	}

	// apply globbing to the ignore index
	if err := globIndex(filesystem, ignoreIndex); err != nil {
		return err
	}

	return fs.WalkDir(filesystem, ".", func(path string, info fs.DirEntry, err error) error {
		if config.shouldIgnore(ignoreIndex, path) {
			if info.IsDir() {
				if debug {
					log.Printf("skip dir: %s", path)
				}

				if statistics {
					defer linter.getStatistic().AddDirSkip()
				}

				return fs.SkipDir
			}

			if debug {
				log.Printf("skip file: %s", path)
			}

			if statistics {
				defer linter.getStatistic().AddFileSkip()
			}

			return nil
		}

		if info == nil {
			return fmt.Errorf("%s not found", path)
		}

		if info.IsDir() {
			if debug {
				log.Printf("lint dir: %s", path)
			}

			if statistics {
				defer linter.getStatistic().AddDir()
			}

			return linter.validateDir(config, index, path)
		}

		if debug {
			log.Printf("lint file: %s", path)
		}

		if statistics {
			defer linter.getStatistic().AddFile()
		}

		return linter.validateFile(config, index, path)
	})
}
