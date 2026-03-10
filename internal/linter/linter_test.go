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
				"snake_case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.jpg":              &fstest.MapFile{Mode: fs.ModePerm},
				"kabab-case.test.jpg":         &fstest.MapFile{Mode: fs.ModePerm},
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
						".*":     "kebab-case",
						".*.jpg": "kebab-case",
						".png":   "snake_case",
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
				Files:     10,
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
		{
			description: "exists with explicit file key",
			filesystem: fstest.MapFS{
				"pkg":           &fstest.MapFile{Mode: fs.ModeDir},
				"pkg/AGENTS.md": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"pkg": config.Ls{
							"AGENTS.md": "exists:1",
						},
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
				Dirs:      2,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists with explicit file key error",
			filesystem: fstest.MapFS{
				"pkg":           &fstest.MapFile{Mode: fs.ModeDir},
				"pkg/README.md": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"pkg": config.Ls{
							"AGENTS.md": "exists:1",
						},
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
				Dirs:      2,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "pkg",
					Ext:  "AGENTS.md",
					Rules: []rule.Rule{
						new(rule.Exists).Init(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "exists with explicit dir key",
			filesystem: fstest.MapFS{
				"packages":         &fstest.MapFile{Mode: fs.ModeDir},
				"packages/app":     &fstest.MapFile{Mode: fs.ModeDir},
				"packages/app/src": &fstest.MapFile{Mode: fs.ModeDir},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "exists",
						},
						"packages/*": config.Ls{
							".dir": "exists:1",
							"src":  "exists:1",
						},
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
				Files:     0,
				FileSkips: 0,
				Dirs:      4,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists with explicit dir key error",
			filesystem: fstest.MapFS{
				"packages":         &fstest.MapFile{Mode: fs.ModeDir},
				"packages/app":     &fstest.MapFile{Mode: fs.ModeDir},
				"packages/app/lib": &fstest.MapFile{Mode: fs.ModeDir},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "exists",
						},
						"packages/*": config.Ls{
							".dir": "exists:1",
							"src":  "exists:1",
						},
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
				Files:     0,
				FileSkips: 0,
				Dirs:      4,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "packages/app",
					Ext:  "src",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "exists supports explicit mandatory keys",
			filesystem: fstest.MapFS{
				"packages":                      &fstest.MapFile{Mode: fs.ModeDir},
				"packages/my-package":           &fstest.MapFile{Mode: fs.ModeDir},
				"packages/my-package/AGENTS.md": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "exists",
						},
						"packages/*": config.Ls{
							".dir":      "kebab-case | exists:1",
							"AGENTS.md": "exists:1",
						},
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
				Dirs:      3,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists explicit file key with error",
			filesystem: fstest.MapFS{
				"packages":                 &fstest.MapFile{Mode: fs.ModeDir},
				"packages/foo":             &fstest.MapFile{Mode: fs.ModeDir},
				"packages/bar":             &fstest.MapFile{Mode: fs.ModeDir},
				"packages/bar/AGENTS.md":   &fstest.MapFile{Mode: fs.ModePerm},
				"packages/bar/another.txt": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "exists",
						},
						"packages/*": config.Ls{
							".dir":      "kebab-case | exists:1",
							"AGENTS.md": "exists:1",
						},
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
				Files:     2,
				FileSkips: 0,
				Dirs:      4,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "packages/foo",
					Ext:  "AGENTS.md",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "exists explicit file key with paths and bypass error",
			filesystem: fstest.MapFS{
				"packages":                     &fstest.MapFile{Mode: fs.ModeDir},
				"packages/foo":                 &fstest.MapFile{Mode: fs.ModeDir},
				"packages/bar":                 &fstest.MapFile{Mode: fs.ModeDir},
				"packages/bar/AGENTS.md":       &fstest.MapFile{Mode: fs.ModePerm},
				"packages/bar/another_file.md": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: map[string]struct{}{
				"packages/bar/AGENTS.md": {},
			},
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "exists",
						},
						"packages/*": config.Ls{
							".dir":      "kebab-case | exists:1",
							"AGENTS.md": "exists:1",
							".md":       "exists:1-2",
						},
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
				Files:     2,
				FileSkips: 0,
				Dirs:      4,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists explicit file key with range combined error",
			filesystem: fstest.MapFS{
				"packages":          &fstest.MapFile{Mode: fs.ModeDir},
				"packages/foo":      &fstest.MapFile{Mode: fs.ModeDir},
				"packages/foo/B.md": &fstest.MapFile{Mode: fs.ModePerm},
				"packages/bar":      &fstest.MapFile{Mode: fs.ModeDir},
				"packages/bar/C.md": &fstest.MapFile{Mode: fs.ModePerm},
				"packages/bar/D.md": &fstest.MapFile{Mode: fs.ModePerm},
				"packages/bar/E.md": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "exists",
						},
						"packages/*": config.Ls{
							".dir":      "kebab-case | exists:1",
							"AGENTS.md": "exists:1",
							".md":       "exists:1-2",
						},
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
				Files:     4,
				FileSkips: 0,
				Dirs:      4,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "packages/bar",
					Ext:  "AGENTS.md",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "packages/bar",
					Ext:  ".md",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1-2"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "packages/foo",
					Ext:  "AGENTS.md",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
			},
		},
		{
			description: "exists monorepo typescript ui example",
			filesystem: fstest.MapFS{
				"packages":                                          &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui":                                       &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/AGENTS.md":                             &fstest.MapFile{Mode: fs.ModePerm},
				"packages/ui/README.md":                             &fstest.MapFile{Mode: fs.ModePerm},
				"packages/ui/CLAUDE.md":                             &fstest.MapFile{Mode: fs.ModePerm},
				"packages/ui/src":                                   &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/src/components":                        &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/src/components/button":                 &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/src/components/button/button.tsx":      &fstest.MapFile{Mode: fs.ModePerm},
				"packages/ui/src/components/button/button.test.tsx": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "kebab-case",
						},
						"packages/*": config.Ls{
							".dir":      "kebab-case",
							".md":       "regex:^(AGENTS|README|CLAUDE|GEMINI)$",
							"AGENTS.md": "exists:1",
							"README.md": "exists:1",
							"src":       "exists:1",
						},
						"packages/ui/src/components": config.Ls{
							".dir": "kebab-case | exists",
							".tsx": "exists:0",
						},
						"packages/ui/src/components/*": config.Ls{
							".tsx":      "regex:${0} | exists:1",
							".test.tsx": "regex:${0} | exists:1",
						},
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
				Files:     5,
				FileSkips: 0,
				Dirs:      6,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{},
		},
		{
			description: "exists monorepo typescript ui example with errors",
			filesystem: fstest.MapFS{
				"packages":                                  &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui":                               &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/AGENTS.md":                     &fstest.MapFile{Mode: fs.ModePerm},
				"packages/ui/src":                           &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/src/components":                &fstest.MapFile{Mode: fs.ModeDir},
				"packages/ui/src/components/Button.tsx":     &fstest.MapFile{Mode: fs.ModePerm},
				"packages/ui/src/components/NOT_ALLOWED.md": &fstest.MapFile{Mode: fs.ModePerm},
			},
			paths: nil,
			linter: NewLinter(
				".",
				config.NewConfig(
					config.Ls{
						"packages": config.Ls{
							".dir": "kebab-case",
						},
						"packages/*": config.Ls{
							".dir":      "kebab-case",
							".md":       "regex:^(AGENTS|README|CLAUDE|GEMINI)$",
							"AGENTS.md": "exists:1",
							"README.md": "exists:1",
							"src":       "exists:1",
						},
						"packages/ui/src/components": config.Ls{
							".dir": "kebab-case | exists",
							".tsx": "exists:0",
						},
						"packages/ui/src/components/*": config.Ls{
							".tsx":      "regex:${0} | exists:1",
							".test.tsx": "regex:${0} | exists:1",
						},
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
				Files:     3,
				FileSkips: 0,
				Dirs:      5,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*rule.Error{
				{
					Path: "packages/ui",
					Ext:  "README.md",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "packages/ui/src/components",
					Ext:  ".tsx",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"0"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "packages/ui/src/components/*",
					Ext:  ".test.tsx",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
					},
					RWMutex: new(sync.RWMutex),
				},
				{
					Path: "packages/ui/src/components/*",
					Ext:  ".tsx",
					Rules: []rule.Rule{
						func() rule.Rule {
							r := new(rule.Exists).Init()
							_ = r.SetParameters([]string{"1"})
							return r
						}(),
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
				expectedRule := test.expectedErrors[j].GetRules()[k]
				if tmpRule.GetName() != expectedRule.GetName() {
					t.Error(equalErrorsErr)
					return
				}

				expectedRuleParameters := expectedRule.GetParameters()
				compareRuleParameters := len(expectedRuleParameters) > 0
				if tmpRule.GetName() == "exists" && reflect.DeepEqual(expectedRuleParameters, new(rule.Exists).Init().GetParameters()) {
					compareRuleParameters = false
				}

				if compareRuleParameters && !reflect.DeepEqual(tmpRule.GetParameters(), expectedRuleParameters) {
					t.Error(equalErrorsErr)
					return
				}
			}
		}

		i++
	}
}

func TestLinter_Run_MonorepoComplexConstraints(t *testing.T) {
	newMonorepoLs := func() config.Ls {
		return config.Ls{
			".dir":                "kebab-case",
			".md":                 "kebab-case | regex:^(README|AGENTS|CLAUDE|GEMINI)$",
			".*":                  "exists:0",
			".json":               "regex:^(package|turbo)$",
			".*.json":             "regex:^tsconfig\\.base$",
			".yaml":               "regex:^pnpm-workspace$",
			"package.json":        "exists:1",
			"pnpm-workspace.yaml": "exists:1",
			"turbo.json":          "exists:0-1",
			"tsconfig.base.json":  "exists:0-1",
			"README.md":           "exists:0-1",
			"AGENTS.md":           "exists:0-1",
			"CLAUDE.md":           "exists:0-1",
			"GEMINI.md":           "exists:0-1",
			"packages": config.Ls{
				".dir": "kebab-case",
			},
			"packages/*": config.Ls{
				".dir":      "kebab-case",
				".md":       "regex:^(AGENTS|README|CLAUDE|GEMINI)$",
				".ts":       "camelCase | PascalCase",
				".tsx":      "camelCase | PascalCase",
				".js":       "camelCase | PascalCase",
				".jsx":      "camelCase | PascalCase",
				"AGENTS.md": "exists:1",
				"README.md": "exists:1",
				"src":       "exists:1",
			},
			"packages/ui/src/components": config.Ls{
				".dir": "kebab-case | exists",
				".tsx": "exists:0",
			},
			"packages/ui/src/components/*": config.Ls{
				".tsx":      "regex:${0} | exists:1",
				".test.tsx": "regex:${0} | exists:1",
			},
		}
	}

	newMonorepoIgnore := func() []string {
		return []string{
			"node_modules",
			".next",
			"coverage",
			"dist",
			"build",
			"packages/ui/dist",
			".env*",
			"**/.env*",
		}
	}

	newMonorepoLinter := func() *Linter {
		return NewLinter(
			".",
			config.NewConfig(
				newMonorepoLs(),
				newMonorepoIgnore(),
			),
			&debug.Statistic{
				Start:     time.Now(),
				Files:     0,
				FileSkips: 0,
				Dirs:      0,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			[]*rule.Error{},
		)
	}

	t.Run("passes with multiple packages and optional whitelisted files", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"package.json":                      &fstest.MapFile{Mode: fs.ModePerm},
			"pnpm-workspace.yaml":               &fstest.MapFile{Mode: fs.ModePerm},
			"turbo.json":                        &fstest.MapFile{Mode: fs.ModePerm},
			"tsconfig.base.json":                &fstest.MapFile{Mode: fs.ModePerm},
			"README.md":                         &fstest.MapFile{Mode: fs.ModePerm},
			"AGENTS.md":                         &fstest.MapFile{Mode: fs.ModePerm},
			"architecture-notes.md":             &fstest.MapFile{Mode: fs.ModePerm},
			".env.local":                        &fstest.MapFile{Mode: fs.ModePerm},
			"node_modules":                      &fstest.MapFile{Mode: fs.ModeDir},
			"node_modules/BAD_NAME.js":          &fstest.MapFile{Mode: fs.ModePerm},
			".next":                             &fstest.MapFile{Mode: fs.ModeDir},
			".next/BAD_NAME.js":                 &fstest.MapFile{Mode: fs.ModePerm},
			"build":                             &fstest.MapFile{Mode: fs.ModeDir},
			"build/BAD_NAME.js":                 &fstest.MapFile{Mode: fs.ModePerm},
			"dist":                              &fstest.MapFile{Mode: fs.ModeDir},
			"dist/BAD_NAME.js":                  &fstest.MapFile{Mode: fs.ModePerm},
			"coverage":                          &fstest.MapFile{Mode: fs.ModeDir},
			"coverage/BAD_NAME.js":              &fstest.MapFile{Mode: fs.ModePerm},
			"packages":                          &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui":                       &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/AGENTS.md":             &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/README.md":             &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/GEMINI.md":             &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/.env.development":      &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/src":                   &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/src/useButton.ts":      &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/src/components":        &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/src/components/button": &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/src/components/button/button.tsx":      &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/src/components/button/button.test.tsx": &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/dist":                     &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/dist/BAD_NAME.tsx":        &fstest.MapFile{Mode: fs.ModePerm},
			"packages/data-model":                  &fstest.MapFile{Mode: fs.ModeDir},
			"packages/data-model/AGENTS.md":        &fstest.MapFile{Mode: fs.ModePerm},
			"packages/data-model/README.md":        &fstest.MapFile{Mode: fs.ModePerm},
			"packages/data-model/src":              &fstest.MapFile{Mode: fs.ModeDir},
			"packages/data-model/src/userModel.ts": &fstest.MapFile{Mode: fs.ModePerm},
		}

		l := newMonorepoLinter()
		err := l.Run(filesystem, nil, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(l.GetErrors()) != 0 {
			t.Fatalf("expected no lint errors, got %+v", l.GetErrors())
		}
	})

	t.Run("fails for required monorepo constraints", func(t *testing.T) {
		filesystem := fstest.MapFS{
			"pnpm-workspace.yaml":                       &fstest.MapFile{Mode: fs.ModePerm},
			"package-lock.json":                         &fstest.MapFile{Mode: fs.ModePerm},
			"tsconfig.app.json":                         &fstest.MapFile{Mode: fs.ModePerm},
			"debug.log":                                 &fstest.MapFile{Mode: fs.ModePerm},
			"NOTES.md":                                  &fstest.MapFile{Mode: fs.ModePerm},
			".env.local":                                &fstest.MapFile{Mode: fs.ModePerm},
			"build":                                     &fstest.MapFile{Mode: fs.ModeDir},
			"build/IGNORED_BAD_NAME.ts":                 &fstest.MapFile{Mode: fs.ModePerm},
			"packages":                                  &fstest.MapFile{Mode: fs.ModeDir},
			"packages/Bad_Pkg":                          &fstest.MapFile{Mode: fs.ModeDir},
			"packages/Bad_Pkg/AGENTS.md":                &fstest.MapFile{Mode: fs.ModePerm},
			"packages/Bad_Pkg/README.md":                &fstest.MapFile{Mode: fs.ModePerm},
			"packages/Bad_Pkg/src":                      &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui":                               &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/AGENTS.md":                     &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/NOTES.md":                      &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/.env.test":                     &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/src":                           &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/src/bad-name.js":               &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/src/components":                &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/src/components/Button.tsx":     &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/src/components/NOT_ALLOWED.md": &fstest.MapFile{Mode: fs.ModePerm},
			"packages/ui/dist":                          &fstest.MapFile{Mode: fs.ModeDir},
			"packages/ui/dist/IGNORED_BAD_NAME.tsx":     &fstest.MapFile{Mode: fs.ModePerm},
		}

		l := newMonorepoLinter()
		err := l.Run(filesystem, nil, true)
		if err != nil {
			t.Fatalf("expected no execution error, got %v", err)
		}
		assertErrorHasRule(t, l.GetErrors(), "", "package.json", "exists")
		assertErrorHasRule(t, l.GetErrors(), "package-lock.json", ".json", "regex")
		assertErrorHasRule(t, l.GetErrors(), "tsconfig.app.json", ".*.json", "regex")
		assertErrorHasRule(t, l.GetErrors(), "", ".*", "exists")
		assertErrorHasRule(t, l.GetErrors(), "NOTES.md", ".md", "kebabcase")
		assertErrorHasRule(t, l.GetErrors(), "packages/Bad_Pkg", ".dir", "kebabcase")
		assertErrorHasRule(t, l.GetErrors(), "packages/ui/NOTES.md", ".md", "regex")
		assertErrorHasRule(t, l.GetErrors(), "packages/ui/src/bad-name.js", ".js", "camelcase")
		assertErrorHasRule(t, l.GetErrors(), "packages/ui", "README.md", "exists")
		assertErrorHasRule(t, l.GetErrors(), "packages/ui/src/components", ".tsx", "exists")
		assertErrorHasRule(t, l.GetErrors(), "packages/ui/src/components/*", ".tsx", "exists")
		assertErrorHasRule(t, l.GetErrors(), "packages/ui/src/components/*", ".test.tsx", "exists")
	})
}

func TestLinter_Run_LargeRepoIgnoreConfigs(t *testing.T) {
	const (
		packageCount                  = 250
		filesPerPackage               = 80
		maxPerformanceRegressionRatio = 5.0
	)

	newLargeRepoLs := func(withRootCatchAll bool) config.Ls {
		ls := config.Ls{
			".dir":                "kebab-case",
			".md":                 "kebab-case | regex:^(README|AGENTS|CLAUDE|GEMINI)$",
			".json":               "regex:^(package|turbo)$",
			".*.json":             "regex:^tsconfig\\.base$",
			".yaml":               "regex:^pnpm-workspace$",
			"package.json":        "exists:1",
			"pnpm-workspace.yaml": "exists:1",
			"turbo.json":          "exists:0-1",
			"tsconfig.base.json":  "exists:0-1",
			"README.md":           "exists:0-1",
			"AGENTS.md":           "exists:0-1",
			"CLAUDE.md":           "exists:0-1",
			"GEMINI.md":           "exists:0-1",
			"packages": config.Ls{
				".dir": "kebab-case",
			},
			"packages/*": config.Ls{
				".dir":      "kebab-case",
				".md":       "regex:^(AGENTS|README|CLAUDE|GEMINI)$",
				".ts":       "camelCase | PascalCase",
				".tsx":      "camelCase | PascalCase",
				".js":       "camelCase | PascalCase",
				".jsx":      "camelCase | PascalCase",
				"AGENTS.md": "exists:1",
				"README.md": "exists:1",
				"src":       "exists:1",
			},
		}
		if withRootCatchAll {
			ls[".*"] = "exists:0"
		}

		return ls
	}

	newLargeRepoIgnore := func() []string {
		return []string{
			"node_modules",
			".next",
			"coverage",
			"dist",
			"build",
			".env*",
			"**/.env*",
		}
	}

	newLargeRepoLinter := func(withRootCatchAll bool) *Linter {
		return NewLinter(
			".",
			config.NewConfig(newLargeRepoLs(withRootCatchAll), newLargeRepoIgnore()),
			&debug.Statistic{
				Start:     time.Now(),
				Files:     0,
				FileSkips: 0,
				Dirs:      0,
				DirSkips:  0,
				RWMutex:   new(sync.RWMutex),
			},
			[]*rule.Error{},
		)
	}

	buildLargeRepoFS := func(packageCount int, filesPerPackage int) fstest.MapFS {
		filesystem := fstest.MapFS{
			"package.json":        &fstest.MapFile{Mode: fs.ModePerm},
			"pnpm-workspace.yaml": &fstest.MapFile{Mode: fs.ModePerm},
			"tsconfig.base.json":  &fstest.MapFile{Mode: fs.ModePerm},
			"README.md":           &fstest.MapFile{Mode: fs.ModePerm},
			"AGENTS.md":           &fstest.MapFile{Mode: fs.ModePerm},
			".env.local":          &fstest.MapFile{Mode: fs.ModePerm},
			"packages":            &fstest.MapFile{Mode: fs.ModeDir},
			"node_modules":        &fstest.MapFile{Mode: fs.ModeDir},
			"node_modules/bad.js": &fstest.MapFile{Mode: fs.ModePerm},
			"dist":                &fstest.MapFile{Mode: fs.ModeDir},
			"dist/bad.js":         &fstest.MapFile{Mode: fs.ModePerm},
		}

		for i := 0; i < packageCount; i++ {
			packageName := fmt.Sprintf("package-%04d", i)
			packageDir := fmt.Sprintf("packages/%s", packageName)
			srcDir := fmt.Sprintf("%s/src", packageDir)

			filesystem[packageDir] = &fstest.MapFile{Mode: fs.ModeDir}
			filesystem[fmt.Sprintf("%s/AGENTS.md", packageDir)] = &fstest.MapFile{Mode: fs.ModePerm}
			filesystem[fmt.Sprintf("%s/README.md", packageDir)] = &fstest.MapFile{Mode: fs.ModePerm}
			filesystem[srcDir] = &fstest.MapFile{Mode: fs.ModeDir}
			filesystem[fmt.Sprintf("%s/.env.test", packageDir)] = &fstest.MapFile{Mode: fs.ModePerm}

			for j := 0; j < filesPerPackage; j++ {
				filesystem[fmt.Sprintf("%s/useFeature%d.ts", srcDir, j)] = &fstest.MapFile{Mode: fs.ModePerm}
			}
		}

		return filesystem
	}

	measureRun := func(t *testing.T, withRootCatchAll bool, filesystem fstest.MapFS) time.Duration {
		t.Helper()

		l := newLargeRepoLinter(withRootCatchAll)
		start := time.Now()
		err := l.Run(filesystem, nil, false)
		duration := time.Since(start)
		if err != nil {
			t.Fatalf("expected no execution error, got %v", err)
		}
		if len(l.GetErrors()) != 0 {
			t.Fatalf("expected no lint errors, got %+v", l.GetErrors())
		}

		return duration
	}

	filesystem := buildLargeRepoFS(packageCount, filesPerPackage)
	t.Logf(
		"large repo simulation contains %d packages with %d source files each (> %d source files total), plus package metadata, ignored env files, and directories",
		packageCount,
		filesPerPackage,
		packageCount*filesPerPackage,
	)

	withRootCatchAll := measureRun(t, true, filesystem)
	withoutRootCatchAll := measureRun(t, false, filesystem)

	t.Logf("large repo run with `.*: exists:0`: %s", withRootCatchAll)
	t.Logf("large repo run without `.*: exists:0`: %s", withoutRootCatchAll)

	if withoutRootCatchAll > 0 {
		ratio := float64(withRootCatchAll) / float64(withoutRootCatchAll)
		t.Logf("large repo performance ratio (with catch-all / without): %.2fx", ratio)
		if ratio > maxPerformanceRegressionRatio {
			t.Fatalf("expected root catch-all config to stay within %.2fx of baseline, got %.2fx", maxPerformanceRegressionRatio, ratio)
		}
	}
}

func assertErrorHasRule(t *testing.T, errors []*rule.Error, path string, ext string, ruleName string) {
	t.Helper()

	for _, lintErr := range errors {
		if lintErr.GetPath() != path || lintErr.GetExt() != ext {
			continue
		}

		for _, tmpRule := range lintErr.GetRules() {
			if tmpRule.GetName() == ruleName {
				return
			}
		}

		t.Fatalf("found error for path=%s ext=%s, but rule=%s was missing: %+v", path, ext, ruleName, lintErr.GetRules())
	}

	t.Fatalf("missing expected lint error path=%s ext=%s rule=%s. actual=%+v", path, ext, ruleName, errors)
}
