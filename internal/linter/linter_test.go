package linter

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/debug"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
)

func TestLinter_Run(t *testing.T) {
	start := time.Now()

	tests := []*struct {
		description       string
		filesystem        fs.FS
		paths             map[string]struct{}
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
			paths: nil,
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
			paths: nil,
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
					Ext:  ".png",
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
			paths: nil,
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
			paths: nil,
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
					Ext:  ".png",
					Rules: []rule.Rule{
						new(rule.SnakeCase).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "glob and glob ignore",
			filesystem: fstest.MapFS{
				"snake_case.png":                  &fstest.MapFile{Mode: fs.ModePerm},
				"src/a/a":                         &fstest.MapFile{Mode: fs.ModeDir},
				"src/a/a/kebab-case.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/a/a/kebab-case.jpg":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/b/b":                         &fstest.MapFile{Mode: fs.ModeDir},
				"src/b/b/kebab-case.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/b/b/kebab-case.jpg":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c":                         &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/PascalCase.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/PascalCase.jpg":          &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/ignore.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/c/packages":                &fstest.MapFile{Mode: fs.ModeDir},
				"src/c/c/packages/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"src/c/d/snake_case.png":          &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case",
						".jpg": "kebab-case",
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
						"src/c/*/*.jpg",
						"src/c/d/*",
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
				Files:     7,
				FileSkips: 3,
				Dirs:      10,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "defaults",
			filesystem: fstest.MapFS{
				// these are valid snake_case file names
				"snake_case.png":     &fstest.MapFile{Mode: fs.ModePerm},
				"snake_case_123.png": &fstest.MapFile{Mode: fs.ModePerm},
				"snake123.png":       &fstest.MapFile{Mode: fs.ModePerm},
				"snake":              &fstest.MapFile{Mode: fs.ModePerm},
				"123.png":            &fstest.MapFile{Mode: fs.ModePerm},

				// this one is caught by the default rule (kebab-case)
				"foobar.tmp": &fstest.MapFile{Mode: fs.ModePerm},

				// these are valid kebab-case file names
				"kebab-case.jpg":      &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case-123.jpg":  &fstest.MapFile{Mode: fs.ModePerm},
				"kebab123.jpg":        &fstest.MapFile{Mode: fs.ModePerm},
				"kebab.jpg":           &fstest.MapFile{Mode: fs.ModePerm},
				"123.jpg":             &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.test.jpg": &fstest.MapFile{Mode: fs.ModePerm},

				// these are valid lowercase file names
				"lowercase.svg":    &fstest.MapFile{Mode: fs.ModePerm},
				"lowercase123.svg": &fstest.MapFile{Mode: fs.ModePerm},
				"123.svg":          &fstest.MapFile{Mode: fs.ModePerm},

				// these are valid flatcase file names
				"flatcase.raw":    &fstest.MapFile{Mode: fs.ModePerm},
				"flatcase123.raw": &fstest.MapFile{Mode: fs.ModePerm},
				"123.raw":         &fstest.MapFile{Mode: fs.ModePerm},

				// these are valid PascalCase file names
				"PascalCase.txt":    &fstest.MapFile{Mode: fs.ModePerm},
				"PascalCase123.txt": &fstest.MapFile{Mode: fs.ModePerm},
				"Pascal.txt":        &fstest.MapFile{Mode: fs.ModePerm},
				"Pascal123.txt":     &fstest.MapFile{Mode: fs.ModePerm},
				"123.txt":           &fstest.MapFile{Mode: fs.ModePerm},

				// these are valid SCREAMING_SNAKE_CASE file names
				"SCREAMING_SNAKE_CASE.bmp": &fstest.MapFile{Mode: fs.ModePerm},
				"SCREAMING123.bmp":         &fstest.MapFile{Mode: fs.ModePerm},
				"123.bmp":                  &fstest.MapFile{Mode: fs.ModePerm},

				"sub":                         &fstest.MapFile{Mode: fs.ModeDir},
				"sub/snake_case.png":          &fstest.MapFile{Mode: fs.ModePerm},
				"sub/kebab-case.jpg":          &fstest.MapFile{Mode: fs.ModePerm},
				"sub/kebab-case.test.jpg":     &fstest.MapFile{Mode: fs.ModePerm},
				"sub/PascalCase.service.jpg":  &fstest.MapFile{Mode: fs.ModePerm},
				"sub/camelCase.service.gif":   &fstest.MapFile{Mode: fs.ModePerm},
				"sub/PascalCase.app.gif":      &fstest.MapFile{Mode: fs.ModePerm},
				"sub/PascalCase.app.test.gif": &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".*":     "kebab-case", // default rule
						".*.jpg": "kebab-case",
						".png":   "snake_case",
						".bmp":   "SCREAMING_SNAKE_CASE",
						".txt":   "PascalCase",
						".svg":   "lowercase",
						".raw":   "flatcase",
						"sub": config.Ls{
							".*":            "kebab-case",
							".*.*":          "kebab-case",
							".service.jpg":  "PascalCase",
							".*.jpg":        "kebab-case",
							".service.*":    "camelCase",
							".app.test.gif": "PascalCase",
							".*.gif":        "PascalCase",
							".png":          "snake_case",
						},
					},
					[]string{
						"node_modules",
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
				Files:     33,
				FileSkips: 0,
				Dirs:      2,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists",
			filesystem: fstest.MapFS{
				"snake_case.png":                     &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.png":                     &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                       &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png":        &fstest.MapFile{Mode: fs.ModePerm},
				"test":                               &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub":                           &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub/snake_case_123.png":        &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/snake_case_456.png":        &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/subsub":                    &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub/subsub/snake_case_123.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/subsub/snake_case_456.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/subsub/service.test.ts":    &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case | exists:1",
						"test": config.Ls{
							".dir": "exists:1",
						},
						"test/*": config.Ls{
							".*":   "exists:0",
							".png": "snake_case | exists:1-2",
							"*": config.Ls{
								".*.ts": "snake_case | exists:1",
							},
						},
						"not_exists": config.Ls{
							".dir": "exists:0",
						},
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
				Files:     6,
				FileSkips: 1,
				Dirs:      4,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists with paths",
			filesystem: fstest.MapFS{
				"snake_case.png":                     &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.png":                     &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                       &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png":        &fstest.MapFile{Mode: fs.ModePerm},
				"test":                               &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub":                           &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub/snake_case_123.png":        &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/snake_case_456.png":        &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/subsub":                    &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub/subsub/snake_case_123.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/subsub/snake_case_456.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: map[string]struct{}{
				"snake_case.png":              {},
				"test":                        {},
				"test/sub/snake_case_123.png": {},
			},
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case | exists:1",
						"test": config.Ls{
							".dir": "exists:1",
						},
						"test/*": config.Ls{
							".*":   "exists:0",
							".png": "snake_case | exists:1-2",
						},
						"not_exists": config.Ls{
							".dir": "exists:0",
						},
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
				Files:     5,
				FileSkips: 1,
				Dirs:      4,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists with error",
			filesystem: fstest.MapFS{
				"snake_case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test":                        &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub":                    &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub/test.ts":            &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/snake_case_123.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/snake_case_456.png": &fstest.MapFile{Mode: fs.ModePerm},
				"wildcards":                   &fstest.MapFile{Mode: fs.ModeDir},
				"wildcards/a":                 &fstest.MapFile{Mode: fs.ModeDir},
				"wildcards/a/b":               &fstest.MapFile{Mode: fs.ModeDir},
				"wildcards/a/b/test.vue":      &fstest.MapFile{Mode: fs.ModePerm},
				"wildcards/a/b/c":             &fstest.MapFile{Mode: fs.ModeDir},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case | exists:2",
						"test": config.Ls{
							".dir": "exists:1",
						},
						"test/*": config.Ls{
							".*":   "exists:0",
							".png": "snake_case | exists:3-5",
						},
						"not_exists": config.Ls{
							".dir": "exists:1",
						},
						"dir_not_exists/*/not_exists": config.Ls{
							".dir": "exists",
						},
						"wildcards/**": config.Ls{
							".dir": "exists:1",
							".*":   "snake_case | exists:1",
							".vue": "snake_case | exists:1",
						},
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
				Files:     5,
				FileSkips: 1,
				Dirs:      7,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "not_exists",
					Ext:  ".dir",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "test/sub",
					Ext:  ".png",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "test/sub",
					Ext:  ".*",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "",
					Ext:  ".png",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards",
					Ext:  ".vue",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards",
					Ext:  ".*",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards/a",
					Ext:  ".vue",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards/a",
					Ext:  ".*",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards/a/b",
					Ext:  ".*",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards/a/b/c",
					Ext:  ".vue",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "wildcards/a/b/c",
					Ext:  ".*",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "dir_not_exists/*/not_exists",
					Ext:  ".dir",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "exists with paths and bypass error",
			filesystem: fstest.MapFS{
				"snake_case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test":                        &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub":                    &fstest.MapFile{Mode: fs.ModeDir},
				"test/sub/test.ts":            &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/snake_case_123.png": &fstest.MapFile{Mode: fs.ModePerm},
				"test/sub/snake_case_456.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: map[string]struct{}{
				"snake_case.png":   {},
				"test/sub/test.ts": {},
			},
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						".png": "snake_case | exists:2",
						"test": config.Ls{
							".dir": "exists:1",
						},
						"test/*": config.Ls{
							".*":   "exists:0",
							".png": "snake_case | exists:3-5",
							".vue": "exists:1",
							"*": config.Ls{
								".dir": "exists:1 | snake_case",
							},
						},
						"not_exists": config.Ls{
							".dir": "exists:1",
						},
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
				Files:     4,
				FileSkips: 1,
				Dirs:      3,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "test/sub",
					Ext:  ".*",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "",
					Ext:  ".png",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
	}

	i := 0
	for _, test := range tests {
		fmt.Printf("Run test %d (%s)\n", i, test.description)

		err := test.linter.Run(test.filesystem, test.paths, true)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("Test %d (%s) failed with unmatched error value - %v", i, test.description, err)
			return
		}

		if !reflect.DeepEqual(test.linter.GetStatistics(), test.expectedStatistic) {
			t.Errorf("Test %d (%s) failed with unmatched linter statistic values\nexpected: %+v\nactual: %+v", i, test.description, test.expectedStatistic, test.linter.GetStatistics())
			return
		}

		equalErrorsErr := fmt.Errorf("Test %d (%s) failed with unmatched linter errors value\nexpected: %+v\nactual: %+v", i, test.description, test.expectedErrors, test.linter.GetErrors())
		if len(test.linter.GetErrors()) != len(test.expectedErrors) {
			t.Error(equalErrorsErr)
			return
		}

		if len(test.linter.GetErrors()) > 0 {
			slices.SortStableFunc(test.linter.GetErrors(), func(a, b *rule.Error) int {
				return cmp.Compare(strings.ToLower(a.GetPath()+a.GetExt()), strings.ToLower(b.GetPath()+b.GetExt()))
			})

			slices.SortStableFunc(test.expectedErrors, func(a, b *rule.Error) int {
				return cmp.Compare(strings.ToLower(a.GetPath()+a.GetExt()), strings.ToLower(b.GetPath()+b.GetExt()))
			})
		}

		var j int
		var tmpError *rule.Error
		for j, tmpError = range test.linter.GetErrors() {
			if tmpError.GetPath() != test.expectedErrors[j].GetPath() {
				t.Error(equalErrorsErr)
				return
			}

			if tmpError.GetExt() != test.expectedErrors[j].GetExt() {
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
