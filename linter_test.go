package main

import (
	"errors"
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
		filesystem        fs.FS
		config            *Config
		linter            *Linter
		expectedErr       error
		expectedStatistic *Statistic
		expectedErrors    []*Error
	}{
		{
			filesystem: fstest.MapFS{
				"snake_case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"kebab-case.png":              &fstest.MapFile{Mode: fs.ModePerm},
				"node_modules":                &fstest.MapFile{Mode: fs.ModeDir},
				"node_modules/snake_case.png": &fstest.MapFile{Mode: fs.ModePerm},
			},
			config: &Config{
				Ls: map[string]interface{}{
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
				Files:     1,
				FileSkips: 1,
				Dirs:      1,
				DirSkips:  1,
				RWMutex:   new(sync.RWMutex),
			},
			expectedErrors: []*Error{},
		},
	}

	var i = 0
	for _, test := range tests {
		err := test.linter.Run(test.filesystem, test.config, true, true)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("Test %d failed with unmatched error value - %v", i, err)
		}

		if !reflect.DeepEqual(test.linter.getStatistic(), test.expectedStatistic) {
			t.Errorf("Test %d failed with unmatched linter statistic values\nexpected: %+v\nactual: %+v", i, test.expectedStatistic, test.linter.getStatistic())
		}

		if !reflect.DeepEqual(test.linter.getErrors(), test.expectedErrors) {
			t.Errorf("Test %d failed with unmatched linter errors value\nexpected: %+v\nactual: %+v", i, test.expectedErrors, test.linter.getErrors())
		}

		i++
	}
}
