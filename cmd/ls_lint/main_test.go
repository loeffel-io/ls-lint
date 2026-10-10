package main

import (
	"maps"
	"path/filepath"
	"testing"
)

func TestNormalizePaths(t *testing.T) {
	workdir := t.TempDir()

	tests := []struct {
		description string
		paths       []string
		expected    map[string]struct{}
		expectedErr bool
	}{
		{description: "no args", paths: nil, expected: nil, expectedErr: false},
		{description: "relative", paths: []string{"src/foo.js"}, expected: map[string]struct{}{"src/foo.js": {}}, expectedErr: false},
		{description: "dot slash", paths: []string{"./src/foo.js"}, expected: map[string]struct{}{"src/foo.js": {}}, expectedErr: false},
		{description: "trailing slash", paths: []string{"src/"}, expected: map[string]struct{}{"src": {}}, expectedErr: false},
		{description: "absolute", paths: []string{filepath.Join(workdir, "src", "foo.js")}, expected: map[string]struct{}{"src/foo.js": {}}, expectedErr: false},
		{description: "stays inside", paths: []string{"src/../lib/a.js"}, expected: map[string]struct{}{"lib/a.js": {}}, expectedErr: false},
		{description: "duplicates", paths: []string{"a.js", "./a.js"}, expected: map[string]struct{}{"a.js": {}}, expectedErr: false},
		{description: "dot", paths: []string{"."}, expected: nil, expectedErr: false},
		{description: "dot with others", paths: []string{"a.js", ".", "b.js"}, expected: nil, expectedErr: false},
		{description: "absolute workdir", paths: []string{workdir}, expected: nil, expectedErr: false},
		{description: "empty", paths: []string{""}, expected: nil, expectedErr: false},
		{description: "parent", paths: []string{"../x.js"}, expected: nil, expectedErr: true},
		{description: "escapes", paths: []string{"src/../../x.js"}, expected: nil, expectedErr: true},
		{description: "absolute outside", paths: []string{filepath.Join(filepath.Dir(workdir), "x.js")}, expected: nil, expectedErr: true},
	}

	i := 0
	for _, test := range tests {
		res, err := normalizePaths(workdir, test.paths)

		if (err != nil) != test.expectedErr {
			t.Errorf("Test %d (%s) failed with unmatched error value - %v", i, test.description, err)
			return
		}

		if (res == nil) != (test.expected == nil) {
			t.Errorf("Test %d (%s) failed with unmatched nil value - %#v", i, test.description, res)
			return
		}

		if !maps.Equal(res, test.expected) {
			t.Errorf("Test %d (%s) failed with unmatched return value - %#v", i, test.description, res)
			return
		}

		i++
	}
}
