package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadExtendsAndMerge(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "node_modules/team/base.yml", `
ls:
  src:
    .js: camelcase
    .ts: camelcase
  shared:
    .ts: camelcase
ignore:
  - node_modules
  - generated
`)
	writeConfig(t, root, "config/vue.yml", `
extends: ../node_modules/team/base.yml
ls:
  src:
    .vue: PascalCase
ignore:
  - dist
`)
	writeConfig(t, root, "config/react.yml", `
ls:
  react:
    .tsx: PascalCase
ignore:
  - dist
`)
	rootConfig := writeConfig(t, root, ".ls-lint.yml", `
extends:
  - config/vue.yml
  - config/react.yml
ls:
  shared:
    .ts: PascalCase
ignore:
  - coverage
`)

	loaded, err := Load([]string{rootConfig})
	if err != nil {
		t.Fatal(err)
	}

	expectedLs := Ls{
		"src": Ls{
			".vue": "PascalCase",
		},
		"shared": Ls{
			".ts": "PascalCase",
		},
		"react": Ls{
			".tsx": "PascalCase",
		},
	}
	if !reflect.DeepEqual(loaded.Ls, expectedLs) {
		t.Errorf("unexpected ls config:\nexpected: %#v\nactual:   %#v", expectedLs, loaded.Ls)
	}

	expectedIgnore := []string{"coverage", "dist", "generated", "node_modules"}
	if !reflect.DeepEqual(loaded.Ignore, expectedIgnore) {
		t.Errorf("unexpected ignore config: expected %v, actual %v", expectedIgnore, loaded.Ignore)
	}
}

func TestLoadAbsolutePathAndExtendsOnlyConfig(t *testing.T) {
	root := t.TempDir()
	base := writeConfig(t, root, "shared/base.yml", `
ls:
  packages:
    .js: kebab-case
`)
	rootConfig := writeConfig(t, root, ".ls-lint.yml", fmt.Sprintf("extends: '%s'\n", strings.ReplaceAll(base, "'", "''")))

	loaded, err := Load([]string{rootConfig})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := loaded.Ls["packages"]; !ok {
		t.Errorf("expected absolute extended config to be loaded: %#v", loaded.Ls)
	}
}

func TestLoadRepeatedRootsPreservesOrder(t *testing.T) {
	root := t.TempDir()
	base := writeConfig(t, root, "base.yml", `
ls:
  src:
    .js: camelcase
  base:
    .md: kebab-case
ignore:
  - generated
`)
	local := writeConfig(t, root, "local.yml", `
ls:
  src:
    .ts: PascalCase
ignore:
  - generated
  - coverage
`)

	loaded, err := Load([]string{base, local})
	if err != nil {
		t.Fatal(err)
	}

	expectedLs := Ls{
		"src": Ls{
			".ts": "PascalCase",
		},
		"base": Ls{
			".md": "kebab-case",
		},
	}
	if !reflect.DeepEqual(loaded.Ls, expectedLs) {
		t.Errorf("unexpected ls config:\nexpected: %#v\nactual:   %#v", expectedLs, loaded.Ls)
	}

	expectedIgnore := []string{"coverage", "generated"}
	if !reflect.DeepEqual(loaded.Ignore, expectedIgnore) {
		t.Errorf("unexpected ignore config: expected %v, actual %v", expectedIgnore, loaded.Ignore)
	}
}

func TestLoadDiamondReappliesSharedConfig(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, "shared.yml", `
ls:
  shared:
    .js: camelcase
`)
	writeConfig(t, root, "left.yml", `
extends: shared.yml
ls:
  shared:
    .js: PascalCase
  left:
    .js: camelcase
`)
	writeConfig(t, root, "right.yml", `
extends: shared.yml
ls:
  right:
    .js: camelcase
`)
	rootConfig := writeConfig(t, root, "root.yml", `
extends:
  - left.yml
  - right.yml
`)

	loaded, err := Load([]string{rootConfig})
	if err != nil {
		t.Fatal(err)
	}

	expectedShared := Ls{".js": "camelcase"}
	if !reflect.DeepEqual(loaded.Ls["shared"], expectedShared) {
		t.Errorf("expected right branch to reapply shared config, got %#v", loaded.Ls["shared"])
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string
		root        string
		expectedErr []string
	}{
		{
			name: "invalid scalar",
			files: map[string]string{
				"root.yml": "extends: 42\n",
			},
			root:        "root.yml",
			expectedErr: []string{"parse config", "extends must be"},
		},
		{
			name: "invalid list item",
			files: map[string]string{
				"root.yml": "extends:\n  - base.yml\n  - 42\n",
			},
			root:        "root.yml",
			expectedErr: []string{"parse config", "extends must be"},
		},
		{
			name: "empty path",
			files: map[string]string{
				"root.yml": "extends: '  '\n",
			},
			root:        "root.yml",
			expectedErr: []string{"parse config", "non-empty"},
		},
		{
			name: "missing child",
			files: map[string]string{
				"root.yml": "extends: missing.yml\n",
			},
			root:        "root.yml",
			expectedErr: []string{`extends "missing.yml"`, "read config"},
		},
		{
			name: "malformed child",
			files: map[string]string{
				"root.yml": "extends: bad.yml\n",
				"bad.yml":  "ls: [\n",
			},
			root:        "root.yml",
			expectedErr: []string{`extends "bad.yml"`, "parse config"},
		},
		{
			name: "self cycle",
			files: map[string]string{
				"root.yml": "extends: root.yml\n",
			},
			root:        "root.yml",
			expectedErr: []string{"config extends cycle", "root.yml"},
		},
		{
			name: "indirect cycle",
			files: map[string]string{
				"a.yml": "extends: b.yml\n",
				"b.yml": "extends: a.yml\n",
			},
			root:        "a.yml",
			expectedErr: []string{"config extends cycle", "a.yml", "b.yml"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			for path, content := range test.files {
				writeConfig(t, root, path, content)
			}

			_, err := Load([]string{filepath.Join(root, test.root)})
			if err == nil {
				t.Fatal("expected load error")
			}

			for _, expected := range test.expectedErr {
				if !strings.Contains(err.Error(), expected) {
					t.Errorf("expected error to contain %q, got %q", expected, err)
				}
			}
		})
	}
}

func writeConfig(t *testing.T, root string, path string, content string) string {
	t.Helper()

	absolutePath := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(absolutePath, []byte(strings.TrimSpace(content)+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	return absolutePath
}
