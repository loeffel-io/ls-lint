package main

import (
	"maps"
	"path/filepath"
	"testing"
)

func TestNormalizePaths(t *testing.T) {
	workdir := t.TempDir()

	tests := []struct {
		name            string
		paths           []string
		normalizedPaths map[string]struct{}
		expectedErr     bool
	}{
		{name: "no args", paths: nil, normalizedPaths: nil},
		{name: "relative", paths: []string{"src/foo.js"}, normalizedPaths: map[string]struct{}{"src/foo.js": {}}},
		{name: "dot slash", paths: []string{"./src/foo.js"}, normalizedPaths: map[string]struct{}{"src/foo.js": {}}},
		{name: "trailing slash", paths: []string{"src/"}, normalizedPaths: map[string]struct{}{"src": {}}},
		{name: "absolute", paths: []string{filepath.Join(workdir, "src", "foo.js")}, normalizedPaths: map[string]struct{}{"src/foo.js": {}}},
		{name: "stays inside", paths: []string{"src/../lib/a.js"}, normalizedPaths: map[string]struct{}{"lib/a.js": {}}},
		{name: "duplicates", paths: []string{"a.js", "./a.js"}, normalizedPaths: map[string]struct{}{"a.js": {}}},
		{name: "dot", paths: []string{"."}, normalizedPaths: nil},
		{name: "dot with others", paths: []string{"a.js", ".", "b.js"}, normalizedPaths: nil},
		{name: "absolute workdir", paths: []string{workdir}, normalizedPaths: nil},
		{name: "empty", paths: []string{""}, normalizedPaths: nil},
		{name: "parent", paths: []string{"../x.js"}, expectedErr: true},
		{name: "escapes", paths: []string{"src/../../x.js"}, expectedErr: true},
		{name: "absolute outside", paths: []string{filepath.Join(filepath.Dir(workdir), "x.js")}, expectedErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizePaths(workdir, test.paths)
			if (err != nil) != test.expectedErr {
				t.Fatalf("err = %v, expectedErr = %v", err, test.expectedErr)
			}

			if (got == nil) != (test.normalizedPaths == nil) || !maps.Equal(got, test.normalizedPaths) {
				t.Fatalf("got = %v, want = %v", got, test.normalizedPaths)
			}
		})
	}
}
