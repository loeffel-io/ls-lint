package glob

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"io/fs"
	"strings"
)

func Index[IndexValue bool | map[string][]rule.Rule](filesystem fs.FS, index map[string]IndexValue, files bool) (err error) {
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
