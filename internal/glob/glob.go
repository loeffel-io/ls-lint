package glob

import (
	"cmp"
	"io/fs"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
)

// Index expands every glob key of the index into the paths it matches. When
// several globs match the same path, the most specific one is applied.
func Index(filesystem fs.FS, index config.RuleIndex, files bool) (err error) {
	var keys []string
	for key := range index {
		if strings.ContainsAny(key, "*{}") {
			keys = append(keys, key)
		}
	}

	slices.SortFunc(keys, compareSpecificity)

	for _, key := range keys {
		var matches []string

		value := index[key]

		if matches, err = doublestar.Glob(filesystem, key); err != nil {
			return err
		}

		if len(matches) == 0 {
			// delete(index, key) // https://github.com/loeffel-io/ls-lint/issues/249
			continue
		}

		for _, match := range matches {
			var matchInfo fs.FileInfo

			if matchInfo, err = fs.Stat(filesystem, match); err != nil {
				return err
			}

			if !files && !matchInfo.IsDir() {
				continue
			}

			if _, ok := index[match]; !ok {
				valueCopy := make(map[string][]rule.Rule, len(value))
				for k, rules := range value {
					valueCopy[k] = make([]rule.Rule, len(rules))
					for i, r := range rules {
						valueCopy[k][i] = r.Copy()
					}
				}

				index[match] = valueCopy
			}

			delete(index, key)
		}
	}

	return nil
}

// compareSpecificity orders glob keys from the most specific to the least:
// more literal path segments first, then fewer wildcards, then by the key
// itself so that the order never depends on map iteration.
func compareSpecificity(a, b string) int {
	return cmp.Or(
		cmp.Compare(literalSegments(b), literalSegments(a)),
		cmp.Compare(strings.Count(a, "*"), strings.Count(b, "*")),
		strings.Compare(a, b),
	)
}

func literalSegments(key string) (count int) {
	for _, segment := range strings.Split(key, "/") {
		if !strings.ContainsAny(segment, "*{}") {
			count++
		}
	}

	return count
}

func IgnoreIndex(filesystem fs.FS, index map[string]bool, files bool) (err error) {
	for key, value := range index {
		var matches []string

		if !strings.ContainsAny(key, "*{}") {
			continue
		}

		if matches, err = doublestar.Glob(filesystem, key); err != nil {
			return err
		}

		if len(matches) == 0 {
			delete(index, key)
			continue
		}

		for _, match := range matches {
			var matchInfo fs.FileInfo

			if matchInfo, err = fs.Stat(filesystem, match); err != nil {
				return err
			}

			if !files && !matchInfo.IsDir() {
				continue
			}

			if _, ok := index[match]; !ok {
				index[match] = value
			}

			delete(index, key)
		}
	}

	return nil
}
