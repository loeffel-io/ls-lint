package main

import (
	"errors"
	"io/fs"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"
)

func TestLinterRun(t *testing.T) {
	var tests = []*struct {
		filesystem     fs.FS
		config         *Config
		linter         *Linter
		expectedErr    error
		expectedErrors []*Error
	}{
		{
			filesystem: fstest.MapFS{
				"snake_case.png": new(fstest.MapFile),
			},
			config: &Config{
				Ls: map[string]interface{}{
					".png": "snake_case",
				},
				Ignore: []string{
					"node_modules",
				},
				RWMutex: new(sync.RWMutex),
			},
			linter: &Linter{
				Errors:  []*Error{},
				RWMutex: new(sync.RWMutex),
			},
			expectedErr:    nil,
			expectedErrors: []*Error{},
		},
	}

	var i = 0
	for _, test := range tests {
		err := test.linter.Run(test.filesystem, test.config)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("Test %d failed with unmatched error value - %v", i, err)
		}

		if !reflect.DeepEqual(test.linter.getErrors(), test.expectedErrors) {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, test.linter.getErrors())
		}

		i++
	}
}
