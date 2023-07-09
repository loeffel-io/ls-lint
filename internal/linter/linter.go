package linter

import (
	"fmt"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	"github.com/loeffel-io/ls-lint/v2/internal/glob"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

const (
	extSep = "."
	dir    = ".dir"
)

type Linter struct {
	root      string
	config    *config.Config
	statistic *debug.Statistic
	errors    []*rule.Error
	*sync.RWMutex
}

func NewLinter(root string, config *config.Config, statistic *debug.Statistic, errors []*rule.Error) *Linter {
	return &Linter{
		root:      root,
		config:    config,
		statistic: statistic,
		errors:    errors,
		RWMutex:   new(sync.RWMutex),
	}
}

func (linter *Linter) GetStatistics() *debug.Statistic {
	linter.RLock()
	defer linter.RUnlock()

	return linter.statistic
}

func (linter *Linter) GetErrors() []*rule.Error {
	linter.RLock()
	defer linter.RUnlock()

	return linter.errors
}

func (linter *Linter) AddError(error *rule.Error) {
	linter.Lock()
	defer linter.Unlock()

	linter.errors = append(linter.errors, error)
}

func (linter *Linter) validateDir(index config.RuleIndex, path string) error {
	var g = new(errgroup.Group)
	var rulesError = 0
	var rulesErrorMutex = new(sync.Mutex)

	rules := linter.config.GetConfig(index, path)
	basename := filepath.Base(path)

	if basename == linter.root {
		return nil
	}

	if _, exists := rules[dir]; !exists {
		return nil
	}

	for _, ruleDir := range rules[dir] {
		ruleDirCopy := ruleDir
		g.Go(func() error {
			valid, err := ruleDirCopy.Validate(basename)

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

	linter.AddError(&rule.Error{
		Path:    path,
		Rules:   rules[dir],
		RWMutex: new(sync.RWMutex),
	})

	return nil
}

func (linter *Linter) validateFile(index config.RuleIndex, path string) error {
	var ext string
	var g = new(errgroup.Group)
	var rulesError = 0
	var rulesErrorMutex = new(sync.Mutex)

	exts := strings.Split(filepath.Base(path), extSep)
	rules := linter.config.GetConfig(index, path)

	for i := 1; i < len(exts); i++ {
		ext = fmt.Sprintf("%s%s", extSep, strings.Join(exts[i:], extSep))
		withoutExt := strings.TrimSuffix(filepath.Base(path), ext)

		if _, exists := rules[ext]; exists {
			for _, ruleFile := range rules[ext] {
				ruleFileCopy := ruleFile
				g.Go(func() error {
					valid, err := ruleFileCopy.Validate(withoutExt)

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

	linter.AddError(&rule.Error{
		Path:    path,
		Rules:   rules[ext],
		RWMutex: new(sync.RWMutex),
	})

	return nil
}

func (linter *Linter) Run(filesystem fs.FS, debug bool, statistics bool) (err error) {
	var index config.RuleIndex
	var ignoreIndex = linter.config.GetIgnoreIndex()

	// create index
	if index, err = linter.config.GetIndex(linter.config.GetLs()); err != nil {
		return err
	}

	// glob index
	if err = glob.Index(filesystem, index, false); err != nil {
		return err
	}

	// glob ignore index
	if err = glob.Index(filesystem, ignoreIndex, true); err != nil {
		return err
	}

	return fs.WalkDir(filesystem, linter.root, func(path string, info fs.DirEntry, err error) error {
		if linter.config.ShouldIgnore(ignoreIndex, path) {
			if info.IsDir() {
				if debug {
					log.Printf("skip dir: %s", path)
				}

				if statistics {
					defer linter.GetStatistics().AddDirSkip()
				}

				return fs.SkipDir
			}

			if debug {
				log.Printf("skip file: %s", path)
			}

			if statistics {
				defer linter.GetStatistics().AddFileSkip()
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
				defer linter.GetStatistics().AddDir()
			}

			return linter.validateDir(index, path)
		}

		if debug {
			log.Printf("lint file: %s", path)
		}

		if statistics {
			defer linter.GetStatistics().AddFile()
		}

		return linter.validateFile(index, path)
	})
}
