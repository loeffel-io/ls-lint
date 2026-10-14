package glob

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/loeffel-io/ls-lint/v2/internal/config"
)

func TestIndex_DeterministicOverlappingGlobs(t *testing.T) {
	filesystem := fstest.MapFS{
		"src/kernel/entities/money.ts": &fstest.MapFile{Mode: fs.ModePerm},
	}

	ls := config.Ls{
		"src/*/*": config.Ls{
			".ts": "kebab-case",
		},
		"src/**/entities": config.Ls{
			".ts": "snake_case",
		},
	}

	var expectedRule string

	for i := range 100 {
		cfg := config.NewConfig(ls, nil)

		index, err := cfg.GetIndex(cfg.GetLs())
		if err != nil {
			t.Fatalf("unexpected GetIndex error on iteration %d: %v", i, err)
		}

		if err = Index(filesystem, index, false); err != nil {
			t.Fatalf("unexpected Index error on iteration %d: %v", i, err)
		}

		rules, ok := index["src/kernel/entities"]
		if !ok {
			t.Fatalf("expected src/kernel/entities in index on iteration %d", i)
		}

		tsRules, ok := rules[".ts"]
		if !ok || len(tsRules) == 0 {
			t.Fatalf("expected .ts rules in src/kernel/entities on iteration %d", i)
		}

		ruleName := tsRules[0].GetName()
		if expectedRule == "" {
			expectedRule = ruleName
		} else if ruleName != expectedRule {
			t.Fatalf("non-deterministic rule on iteration %d: got %q, want %q", i, ruleName, expectedRule)
		}
	}
}

func TestIndex_SingleGlob(t *testing.T) {
	filesystem := fstest.MapFS{
		"packages/app/index.ts": &fstest.MapFile{Mode: fs.ModePerm},
	}

	ls := config.Ls{
		"packages/*": config.Ls{
			".ts": "kebab-case",
		},
	}

	cfg := config.NewConfig(ls, nil)
	index, err := cfg.GetIndex(cfg.GetLs())
	if err != nil {
		t.Fatalf("unexpected GetIndex error: %v", err)
	}

	if err = Index(filesystem, index, false); err != nil {
		t.Fatalf("unexpected Index error: %v", err)
	}

	rules, ok := index["packages/app"]
	if !ok {
		t.Fatal("expected packages/app in index")
	}

	tsRules, ok := rules[".ts"]
	if !ok || len(tsRules) == 0 {
		t.Fatal("expected .ts rules in packages/app")
	}

	if tsRules[0].GetName() != "kebabcase" {
		t.Fatalf("got rule %q, want kebabcase", tsRules[0].GetName())
	}
}

func TestIndex_NoMatches(t *testing.T) {
	filesystem := fstest.MapFS{
		"src/index.ts": &fstest.MapFile{Mode: fs.ModePerm},
	}

	ls := config.Ls{
		"nonexistent/*": config.Ls{
			".ts": "kebab-case",
		},
	}

	cfg := config.NewConfig(ls, nil)
	index, err := cfg.GetIndex(cfg.GetLs())
	if err != nil {
		t.Fatalf("unexpected GetIndex error: %v", err)
	}

	if err = Index(filesystem, index, false); err != nil {
		t.Fatalf("unexpected Index error: %v", err)
	}

	if _, ok := index["nonexistent/*"]; !ok {
		t.Fatal("expected unmatched glob key to remain in index")
	}
}

func TestIgnoreIndex_DeterministicOverlappingGlobs(t *testing.T) {
	filesystem := fstest.MapFS{
		"src/kernel/entities/money.ts": &fstest.MapFile{Mode: fs.ModePerm},
	}

	for i := range 100 {
		ignoreIndex := map[string]bool{
			"src/*/*":         true,
			"src/**/entities": true,
		}

		if err := IgnoreIndex(filesystem, ignoreIndex, false); err != nil {
			t.Fatalf("unexpected IgnoreIndex error on iteration %d: %v", i, err)
		}

		if !ignoreIndex["src/kernel/entities"] {
			t.Fatalf("expected src/kernel/entities ignored on iteration %d", i)
		}
	}
}

func TestIgnoreIndex_NoMatches(t *testing.T) {
	filesystem := fstest.MapFS{
		"src/index.ts": &fstest.MapFile{Mode: fs.ModePerm},
	}

	ignoreIndex := map[string]bool{
		"nonexistent/*": true,
	}

	if err := IgnoreIndex(filesystem, ignoreIndex, false); err != nil {
		t.Fatalf("unexpected IgnoreIndex error: %v", err)
	}

	if _, ok := ignoreIndex["nonexistent/*"]; ok {
		t.Fatal("expected unmatched glob key to be removed from ignore index")
	}
}
