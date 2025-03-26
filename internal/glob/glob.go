package glob

import (
	"io/fs"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/loeffel-io/ls-lint/v2/internal/config"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
)

func Index(filesystem fs.FS, index config.RuleIndex, files bool) (err error) {
	for key, value := range index {
		var matches []string

		if !strings.ContainsAny(key, "*{}") {
			continue
		}

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
