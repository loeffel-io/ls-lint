package main

import (
	"bytes"
	"testing"
)

// TestNormalizeConfig tests normalizeConfig
// includes nonsense regex patterns for testing
func TestNormalizeConfig(t *testing.T) {
	const runeWindowsSep = '\\'

	var tests = []struct {
		unixSep  byte
		sep      byte
		bytes    []byte
		expected []byte
	}{
		{
			unixSep: byte(runeUnixSep),
			sep:     byte(runeUnixSep),
			bytes: []byte(
				`ls:
				  .dir: kebab-case
				  .go: snake_case
				
				  path/to/test:
					.go: snake_case
					.js: kebab-case | camelCase
				
				ignore:
				  - .idea
				  - .git`,
			),
			expected: []byte(
				`ls:
				  .dir: kebab-case
				  .go: snake_case
				
				  path/to/test:
					.go: snake_case
					.js: kebab-case | camelCase
				
				ignore:
				  - .idea
				  - .git`,
			),
		},
		{
			unixSep: byte(runeUnixSep),
			sep:     byte(runeWindowsSep),
			bytes: []byte(
				`ls:
				  .dir: kebab-case
				  .go: snake_case
				
				  path/to/test:
					.go: snake_case
					.js: kebab-case | camelCase
					.py: regex:^[/\-\\]+$
				
				ignore:
				  - path/to/ignore
				  - .idea
				  - .git`,
			),
			expected: []byte(
				`ls:
				  .dir: kebab-case
				  .go: snake_case
				
				  path\to\test:
					.go: snake_case
					.js: kebab-case | camelCase
					.py: regex:^[/\-\\]+$
				
				ignore:
				  - path\to\ignore
				  - .idea
				  - .git`,
			),
		},
		{
			unixSep: byte(runeWindowsSep),
			sep:     byte(runeWindowsSep),
			bytes: []byte(
				`ls:
				  .dir: kebab-case
				  .go: snake_case
				
				  path\to\test:
					.go: snake_case
					.js: kebab-case | camelCase
					.py: regex:^[/\-\\]+$
				
				ignore:
				  - path\to\ignore
				  - .idea
				  - .git`,
			),
			expected: []byte(
				`ls:
				  .dir: kebab-case
				  .go: snake_case
				
				  path\to\test:
					.go: snake_case
					.js: kebab-case | camelCase
					.py: regex:^[/\-\\]+$
				
				ignore:
				  - path\to\ignore
				  - .idea
				  - .git`,
			),
		},
	}

	var i = 0
	for _, test := range tests {
		res := normalizeConfig(test.bytes, test.unixSep, test.sep)

		if !bytes.Equal(res, test.expected) {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}
