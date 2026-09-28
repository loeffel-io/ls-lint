package glob

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/loeffel-io/ls-lint/v2/internal/config"
)

func TestIndex_OverlappingGlobs(t *testing.T) {
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

	// Map iteration order is random, so a single run could pass by chance.
	for range 100 {
		cfg := config.NewConfig(ls, nil)

		index, err := cfg.GetIndex(cfg.GetLs())
		if err != nil {
			t.Fatal(err)
		}

		if err = Index(filesystem, index, false); err != nil {
			t.Fatal(err)
		}

		if got := index["src/kernel/entities"][".ts"][0].GetName(); got != "snakecase" {
			t.Fatalf("src/kernel/entities got %s rules, want the more specific src/**/entities (snakecase)", got)
		}
	}
}
