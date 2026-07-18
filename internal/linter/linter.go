package linter

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	"github.com/loeffel-io/ls-lint/v2/internal/glob"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"golang.org/x/sync/errgroup"
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

func (linter *Linter) validateDir(index config.RuleIndex, path string, validate bool) (string, string, error) {
	indexDir, rules := linter.config.GetConfig(index, path)

	if !validate {
		return indexDir, dir, nil
	}

	g := new(errgroup.Group)

	var rulesNonExclusiveCount int8
	var rulesNonExclusiveError int8
	rulesMutex := new(sync.Mutex)

	var pathDir string
	if pathDir = path; pathDir == "." {
		pathDir = ""
	}

	basename := filepath.Base(path)

	if basename == linter.root {
		return indexDir, dir, nil
	}

	if _, exists := rules[dir]; !exists {
		return indexDir, dir, nil
	}

	for _, ruleDir := range rules[dir] {
		g.Go(func() error {
			if ruleDir.GetName() == "exists" && pathDir != indexDir {
				return nil
			}

			valid, err := ruleDir.Validate(basename, pathDir, ruleDir.GetName() != "exists")
			if err != nil {
				return err
			}

			if !ruleDir.GetExclusive() {
				rulesMutex.Lock()
				rulesNonExclusiveCount++
				if !valid {
					rulesNonExclusiveError++
				}
				rulesMutex.Unlock()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return indexDir, dir, err
	}

	if rulesNonExclusiveError == 0 || rulesNonExclusiveError != rulesNonExclusiveCount {
		return indexDir, dir, nil
	}

	linter.AddError(&rule.Error{
		Path:    path,
		Ext:     dir,
		Rules:   rules[dir],
		RWMutex: new(sync.RWMutex),
	})

	return indexDir, dir, nil
}

func (linter *Linter) validateFile(index config.RuleIndex, path string, validate bool) (string, string, error) {
	var ext string
	g := new(errgroup.Group)

	var rulesNonExclusiveCount int8
	var rulesNonExclusiveError int8
	rulesMutex := new(sync.Mutex)

	indexDir, rules := linter.config.GetConfig(index, path)

	var pathDir string
	pathDir = filepath.ToSlash(filepath.Dir(path)); // compatibility with windows
	if pathDir == "." {
		pathDir = ""
	}

	basename := filepath.Base(path)
	parts := strings.Split(basename, extSep)
	var withoutExt string

	runValidation := func() {
		for _, ruleFile := range rules[ext] {
			if !validate && ruleFile.GetName() != "exists" {
				continue
			}

			g.Go(func() error {
				if ruleFile.GetName() == "exists" && pathDir != indexDir {
					return nil
				}

				valid, err := ruleFile.Validate(withoutExt, pathDir, ruleFile.GetName() != "exists")
				if err != nil {
					return err
				}

				if !ruleFile.GetExclusive() {
					rulesMutex.Lock()
					rulesNonExclusiveCount++
					if !valid {
						rulesNonExclusiveError++
					}
					rulesMutex.Unlock()
				}

				return nil
			})
		}
	}

	if len(parts) == 1 {
		// File without extension
		withoutExt = basename
		potentialExts := []string{basename + extSep, extSep, extSep + "*"}

		for _, currentExt := range potentialExts {
			if _, ok := rules[currentExt]; ok {
				ext = currentExt
				runValidation()
				break
			}
		}
	} else {
		// File with extension(s)
		exts := parts[1:]
		n := len(exts)
		maxCombinations := 1 << n // 2^n combinations

		originalExt := extSep + strings.Join(exts, extSep)
		withoutExt = strings.TrimSuffix(basename, originalExt)

		for i := 0; i < maxCombinations; i++ {
			combination := make([]string, n)
			for j := 0; j < n; j++ {
				if i&(1<<(n-1-j)) == 0 {
					combination[j] = exts[j]
				} else {
					combination[j] = "*"
				}
			}

			currentExt := extSep + strings.Join(combination, extSep)
			if _, ok := rules[currentExt]; ok {
				ext = currentExt
				runValidation()
				break
			}
		}
	}

	if err := g.Wait(); err != nil {
		return indexDir, ext, err
	}

	if !validate || rulesNonExclusiveError == 0 || rulesNonExclusiveError != rulesNonExclusiveCount {
		return indexDir, ext, nil
	}

	linter.AddError(&rule.Error{
		Path:    path,
		Dir:     false,
		Ext:     ext,
		Rules:   rules[ext],
		RWMutex: new(sync.RWMutex),
	})

	return indexDir, ext, nil
}

func (linter *Linter) Run(filesystem fs.FS, paths map[string]struct{}, debug bool) (err error) {
	var pathsIndex map[string]map[string]struct{} = nil
	if len(paths) > 0 {
		pathsIndex = make(map[string]map[string]struct{})
	}

	// create index
	var index config.RuleIndex
	if index, err = linter.config.GetIndex(linter.config.GetLs()); err != nil {
		return err
	}

	// glob index
	if err = glob.Index(filesystem, index, false); err != nil {
		return err
	}

	// glob ignore index
	ignoreIndex := linter.config.GetIgnoreIndex()
	if err = glob.IgnoreIndex(filesystem, ignoreIndex, true); err != nil {
		return err
	}

	if debug {
		fmt.Printf("=============================\nls index\n-----------------------------\n")
		for path, pathIndex := range index {
			switch path == "" {
			case true:
				fmt.Printf(".:")
			case false:
				fmt.Printf("%s:", path)
			}

			for ext, rules := range pathIndex {
				tmpRules := make([]string, 0)
				for _, tmpRule := range rules {
					if len(tmpRule.GetParameters()) > 0 {
						tmpRules = append(tmpRules, fmt.Sprintf("%s:%s", tmpRule.GetName(), strings.Join(tmpRule.GetParameters(), ",")))
						continue
					}

					tmpRules = append(tmpRules, tmpRule.GetName())
				}

				fmt.Printf(" %s: %s", ext, strings.Join(tmpRules, ", "))
			}
			fmt.Printf("\n")
		}

		fmt.Printf("-----------------------------\nignore index\n-----------------------------\n")
		for path := range ignoreIndex {
			fmt.Printf("%s\n", path)
		}

		fmt.Printf("-----------------------------\nlint\n-----------------------------\n")
	}

	if debug {
		defer func() {
			fmt.Printf("-----------------------------\nstatistics\n-----------------------------\n")
			fmt.Printf("time: %s\n", time.Since(linter.GetStatistics().Start).Truncate(time.Microsecond).String())
			fmt.Printf("paths: %d\n", linter.GetStatistics().Files)
			fmt.Printf("file skips: %d\n", linter.GetStatistics().FileSkips)
			fmt.Printf("dirs: %d\n", linter.GetStatistics().Dirs)
			fmt.Printf("dir skips: %d\n", linter.GetStatistics().DirSkips)
			fmt.Printf("=============================\n")
		}()
	}

	if err = fs.WalkDir(filesystem, linter.root, func(path string, info fs.DirEntry, err error) error {
		if linter.config.ShouldIgnore(ignoreIndex, path) {
			if info.IsDir() {
				if debug {
					fmt.Printf("skip dir: %s\n", path)
					linter.GetStatistics().AddDirSkip()
				}

				return fs.SkipDir
			}

			if debug {
				fmt.Printf("skip file: %s\n", path)
				linter.GetStatistics().AddFileSkip()
			}

			return nil
		}

		if info == nil {
			return fmt.Errorf("%s not found", path)
		}

		var indexDir, ext string
		validate := len(paths) == 0
		if _, ok := paths[path]; !validate {
			validate = ok
		}

		if info.IsDir() {
			if debug {
				fmt.Printf("lint dir: %s\n", path)
				linter.GetStatistics().AddDir()
			}

			if indexDir, ext, err = linter.validateDir(index, path, validate); err != nil {
				return err
			}

			if pathsIndex != nil && validate {
				if _, ok := pathsIndex[indexDir]; !ok {
					pathsIndex[indexDir] = make(map[string]struct{})
				}

				pathsIndex[indexDir][ext] = struct{}{}
			}

			return nil
		}

		if debug {
			fmt.Printf("lint file: %s\n", path)
			linter.GetStatistics().AddFile()
		}

		if indexDir, ext, err = linter.validateFile(index, path, validate); err != nil {
			return err
		}

		if pathsIndex != nil && validate {
			if _, ok := pathsIndex[indexDir]; !ok {
				pathsIndex[indexDir] = make(map[string]struct{})
			}

			pathsIndex[indexDir][ext] = struct{}{}
		}

		return nil
	}); err != nil {
		return err
	}

	// validate exists
	for path, pathIndex := range index {
		for ext, rules := range pathIndex {
			if _, ok := pathsIndex[path][ext]; pathsIndex != nil && !ok {
				continue
			}

			for _, r := range rules {
				if r.GetName() != "exists" {
					continue
				}

				var valid bool
				if valid, err = r.Validate("", "", true); err != nil {
					return err
				}

				if !valid {
					linter.AddError(&rule.Error{
						Path:    path,
						Dir:     true,
						Ext:     ext,
						Rules:   []rule.Rule{r},
						RWMutex: new(sync.RWMutex),
					})
				}
			}
		}
	}

	return nil
}
