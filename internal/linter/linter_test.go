package linter

import (
	"errors"
	"fmt"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
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
		linter            *Linter
		expectedErr       error
		expectedStatistic *debug.Statistic
		expectedErrors    []*rule.Error
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
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case",
					},
					[]string{
						"node_modules",
						"kebab-case.png",
					},
				),
				&debug.Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				[]*rule.Error{},
			),
			expectedErr: nil,
			expectedStatistic: &debug.Statistic{
				Start:     start,
				Files:     2,
				FileSkips: 1,
				Dirs:      2,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "fail",
			filesystem: fstest.MapFS{
				"not-snake-case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case",
					},
					[]string{},
				),
				&debug.Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				[]*rule.Error{},
			),
			expectedErr: nil,
			expectedStatistic: &debug.Statistic{
				Start:     start,
				Files:     1,
				FileSkips: 0,
				Dirs:      1,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "not-snake-case.png",
					Rules: []rule.Rule{
						new(rule.SnakeCase).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "glob",
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
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case",
						"src/**/c": config.Ls{
							".png": "PascalCase",
							"packages": config.Ls{
								".png": "snake_case",
							},
						},
						"src/{a,b}/*": config.Ls{
							".png": "kebab-case",
						},
					},
					[]string{
						"src/c/c/ignore.png",
					},
				),
				&debug.Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				[]*rule.Error{},
			),
			expectedErr: nil,
			expectedStatistic: &debug.Statistic{
				Start:     start,
				Files:     5,
				FileSkips: 1,
				Dirs:      9,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "glob (fail)",
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
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case",
						"src/**/c": config.Ls{
							".png": "PascalCase",
							"packages": config.Ls{
								".png": "snake_case",
							},
						},
						"src/{a,b}/*": config.Ls{
							".png": "kebab-case",
						},
					},
					[]string{
						"src/c/c/ignore.png",
					},
				),
				&debug.Statistic{
					Start:     start,
					Files:     0,
					FileSkips: 0,
					Dirs:      0,
					DirSkips:  0,
					RWMutex:   new(sync.RWMutex),
				},
				[]*rule.Error{},
			),
			expectedErr: nil,
			expectedStatistic: &debug.Statistic{
				Start:     start,
				Files:     5,
				FileSkips: 1,
				Dirs:      9,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "src/c/c/packages/not-snake-case.png",
					Rules: []rule.Rule{
						new(rule.SnakeCase).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
	}

	var i = 0
	for _, test := range tests {
		fmt.Printf("Run test %d (%s)\n", i, test.description)

		var err = test.linter.Run(test.filesystem, true, true)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("Test %d (%s) failed with unmatched error value - %v", i, test.description, err)
			return
		}

		if !reflect.DeepEqual(test.linter.GetStatistics(), test.expectedStatistic) {
			t.Errorf("Test %d (%s) failed with unmatched linter statistic values\nexpected: %+v\nactual: %+v", i, test.description, test.expectedStatistic, test.linter.GetStatistics())
			return
		}

		var equalErrorsErr = fmt.Errorf("Test %d (%s) failed with unmatched linter errors value\nexpected: %+v\nactual: %+v", i, test.description, test.expectedErrors, test.linter.GetErrors())
		if len(test.linter.GetErrors()) != len(test.expectedErrors) {
			t.Error(equalErrorsErr)
			return
		}

		var j int
		var tmpError *rule.Error
		for j, tmpError = range test.linter.GetErrors() {
			if tmpError.GetPath() != test.expectedErrors[j].GetPath() {
				t.Error(equalErrorsErr)
				return
			}

			if len(tmpError.GetRules()) != len(test.expectedErrors[j].GetRules()) {
				t.Error(equalErrorsErr)
				return
			}

			var k int
			var tmpRule rule.Rule
			for k, tmpRule = range tmpError.GetRules() {
				if tmpRule.GetName() != test.expectedErrors[j].GetRules()[k].GetName() {
					t.Error(equalErrorsErr)
					return
				}
			}
		}

		i++
	}
}
