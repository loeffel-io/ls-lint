package main

import (
	"errors"
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"
	"time"
)

func emptyStatistics(start time.Time) *Statistic {
	return &Statistic{
		Start:     start,
		Files:     0,
		FileSkips: 0,
		Dirs:      0,
		DirSkips:  0,
		RWMutex:   new(sync.RWMutex),
	}
}

func TestLinterRun(t *testing.T) {
	var start = time.Now()

	var tests = []*struct {
		description       string
		filesystem        fs.FS
		config            *Config
		linter            *Linter
		expectedErr       error
		expectedStatistic *Statistic
		expectedErrors    []*Error
	}{
		{
			description: "No violations detected",
			filesystem: fstest.MapFS{
				"snake_case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test":                        &fstest.MapFile{Mode: fs.ModeDir},
				"test/snake_case_123.png":     &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
				},
				Ignore: []string{
					"node_modules",
					"kebab-case.png",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     2,
				FileSkips: 1,
				Dirs:      2,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
		},
		{
			description: "Single file violation",
			filesystem: fstest.MapFS{
				"not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
				},
				Ignore:  []string{},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     1,
				FileSkips: 0,
				Dirs:      1,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{
				{
					Path: "not-snake-case.png",
					Rules: []Rule{
						new(RuleSnakeCase).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "No violations with globs in config",
			filesystem: fstest.MapFS{
				"snake_case.png":                  &fstest.MapFile{Mode: fs.ModePerm},
				"src/a/a":                         &fstest.MapFile{Mode: fs.ModeDir},
				"src/a/a/kebab-case.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/b/b":                         &fstest.MapFile{Mode: fs.ModeDir},
				"src/b/b/kebab-case.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c":                         &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/PascalCase.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/ignore.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/packages":                &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/packages/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
					"src/**/c": ls{
						".png": "PascalCase",
						"packages": ls{
							".png": "snake_case",
						},
					},
					"src/{a,b}/*": ls{
						".png": "kebab-case",
					},
				},
				Ignore: []string{
					"src/c/c/ignore.png",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     5,
				FileSkips: 1,
				Dirs:      9,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
		},
		{
			description: "Violations with glob in config",
			filesystem: fstest.MapFS{
				"snake_case.png":                      &fstest.MapFile{Mode: fs.ModePerm},
				"src/a/a":                             &fstest.MapFile{Mode: fs.ModeDir},
				"src/a/a/kebab-case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"src/b/b":                             &fstest.MapFile{Mode: fs.ModeDir},
				"src/b/b/kebab-case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c":                             &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/PascalCase.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/ignore.png":                  &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/packages":                    &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/packages/not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
					"src/**/c": ls{
						".png": "PascalCase",
						"packages": ls{
							".png": "snake_case",
						},
					},
					"src/{a,b}/*": ls{
						".png": "kebab-case",
					},
				},
				Ignore: []string{
					"src/c/c/ignore.png",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     5,
				FileSkips: 1,
				Dirs:      9,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{
				{
					Path: "src/c/c/packages/not-snake-case.png",
					Rules: []Rule{
						new(RuleSnakeCase).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "Invalid glob in config",
			filesystem:  fstest.MapFS{},
			config: &Config{
				Ls: ls{
					"src/{a,b/*": ls{
						".png": "kebab-case",
					},
				},
				Ignore:  []string{},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr:       doublestar.ErrBadPattern,
			expectedStatistic: emptyStatistics(start),
			expectedErrors:    []*Error{},
		},
		{
			description: "No violations with glob in ignores",
			filesystem: fstest.MapFS{
				"a/c":                         &fstest.MapFile{Mode: fs.ModeDir},
				"b/c":                         &fstest.MapFile{Mode: fs.ModeDir},
				"a/c/not-snake-case.png":      &fstest.MapFile{Mode: fs.ModePerm},
				"a/c/also-not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"b/c/not-snake-case.png":      &fstest.MapFile{Mode: fs.ModePerm},
				"b/c/also-not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
				},
				Ignore: []string{
					"*/c",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     0,
				FileSkips: 0,
				Dirs:      3,
				DirSkips:  2,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
		},
		{
			description: "No violations with alternatives in ignores",
			filesystem: fstest.MapFS{
				"a/c":                    &fstest.MapFile{Mode: fs.ModeDir},
				"b/c":                    &fstest.MapFile{Mode: fs.ModeDir},
				"a/c/not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"b/c/not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
				},
				Ignore: []string{
					"{a,b}/c",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     0,
				FileSkips: 0,
				Dirs:      3,
				DirSkips:  2,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
		},
		{
			description: "Invalid glob in ignore",
			filesystem:  fstest.MapFS{},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
				},
				Ignore: []string{
					"{a/c",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Statistic: emptyStatistics(start),
				Errors:    []*Error{},
				RWMutex:   new(sync.RWMutex),
			},
			expectedErr:       doublestar.ErrBadPattern,
			expectedStatistic: emptyStatistics(start),
			expectedErrors:    []*Error{},
		},
	}

	var i = 0
	for _, test := range tests {
		fmt.Printf("Run test %d (%s)\n", i, test.description)

		err := test.linter.Run(test.filesystem, test.config, true, true)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("Test %d (%s) failed with unmatched error value - %v", i, test.description, err)
			return
		}

		if !reflect.DeepEqual(test.linter.getStatistic(), test.expectedStatistic) {
			t.Errorf("Test %d (%s) failed with unmatched linter statistic values\nexpected: %+v\nactual: %+v",
				i, test.description, test.expectedStatistic, test.linter.getStatistic())
			return
		}

		var equalErrorsErr = fmt.Errorf("Test %d (%s) failed with unmatched linter errors value\nexpected: %+v\nactual: %+v",
			i, test.description, test.expectedErrors, test.linter.getErrors())

		if len(test.linter.getErrors()) != len(test.expectedErrors) {
			t.Error(equalErrorsErr)
			return
		}

		for i, tmpError := range test.linter.getErrors() {
			if tmpError.getPath() != test.expectedErrors[i].getPath() {
				t.Error(equalErrorsErr)
				return
			}

			if len(tmpError.getRules()) != len(test.expectedErrors[i].getRules()) {
				t.Error(equalErrorsErr)
				return
			}

			for j, tmpRule := range tmpError.getRules() {
				if tmpRule.GetName() != test.expectedErrors[i].getRules()[j].GetName() {
					t.Error(equalErrorsErr)
					return
				}
			}
		}

		i++
	}
}
