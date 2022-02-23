package main

import (
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"
	"time"
)

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
			description: "success",
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
				Statistic: &Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				Errors:  []*Error{},
				RWMutex: new(sync.RWMutex),
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
			description: "fail",
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
				Statistic: &Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				Errors:  []*Error{},
				RWMutex: new(sync.RWMutex),
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
			description: "glob",
			filesystem: fstest.MapFS{
				"snake_case.png":         &fstest.MapFile{Mode: fs.ModePerm},
				"src/a/a":                &fstest.MapFile{Mode: fs.ModeDir},
				"src/a/a/kebab-case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"src/b/b":                &fstest.MapFile{Mode: fs.ModeDir},
				"src/b/b/kebab-case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c":                &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/PascalCase.png": &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/ignore.png":     &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: ls{
					".png": "snake_case",
					"src/**/c": ls{
						".png": "PascalCase",
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
				Statistic: &Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				Errors:  []*Error{},
				RWMutex: new(sync.RWMutex),
			},
			expectedErr: nil,
			expectedStatistic: &Statistic{
				Start:     start,
				Files:     4,
				FileSkips: 1,
				Dirs:      8,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
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
			t.Errorf("Test %d (%s) failed with unmatched linter statistic values\nexpected: %+v\nactual: %+v", i, test.description, test.expectedStatistic, test.linter.getStatistic())
			return
		}

		var equalErrorsErr = fmt.Errorf("Test %d (%s) failed with unmatched linter errors value\nexpected: %+v\nactual: %+v", i, test.description, test.expectedErrors, test.linter.getErrors())
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
