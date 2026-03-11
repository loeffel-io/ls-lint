package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/loeffel-io/ls-lint/v2/internal/rule"
)

func TestGetConfig(t *testing.T) {
	config := new(Config)
	indexMock := map[string]map[string][]rule.Rule{
		".": {
			".dir": []rule.Rule{rule.RulesIndex["lowercase"]},
		},
		"./src": {
			".dir": []rule.Rule{rule.RulesIndex["camelcase"]},
		},
	}
	indexMockEmpty := make(RuleIndex)

	tests := []*struct {
		config   *Config
		index    RuleIndex
		path     string
		expected map[string][]rule.Rule
	}{
		{
			config: config,
			index:  indexMock,
			path:   "./src/test/Test.js",
			expected: map[string][]rule.Rule{
				".dir": {rule.RulesIndex["camelcase"]},
			},
		},
		{
			config: config,
			index:  indexMock,
			path:   "./images/path.png",
			expected: map[string][]rule.Rule{
				".dir": {rule.RulesIndex["lowercase"]},
			},
		},
		{
			config:   config,
			index:    indexMockEmpty,
			path:     "./images/path.png",
			expected: nil,
		},
	}

	i := 0
	for _, test := range tests {
		_, res := test.config.GetConfig(test.index, test.path)

		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}

func TestGetIgnoreIndex(t *testing.T) {
	tests := []struct {
		description   string
		config        *Config
		expectedExact map[string]bool
		expectedGlob  []string
		expectedErr   string
	}{
		{
			description: "splits exact and glob ignores",
			config: NewConfig(nil, []string{
				"node_modules",
				".env*",
				"packages/*/dist",
				`literal\*name`,
			}),
			expectedExact: map[string]bool{
				"node_modules":  true,
				`literal\*name`: true,
			},
			expectedGlob: []string{
				".env*",
				"packages/*/dist",
			},
		},
		{
			description: "fails for invalid glob ignore",
			config: NewConfig(nil, []string{
				"[",
			}),
			expectedErr: `invalid ignore pattern "["`,
		},
	}

	for _, test := range tests {
		index, err := test.config.GetIgnoreIndex()
		if test.expectedErr != "" {
			if err == nil || !errors.Is(err, ErrInvalidIgnorePattern) {
				t.Fatalf("%s: expected invalid ignore pattern error %q, got %v", test.description, test.expectedErr, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: expected no error, got %v", test.description, err)
		}
		if !reflect.DeepEqual(index.Exact, test.expectedExact) {
			t.Fatalf("%s: expected exact index %+v, got %+v", test.description, test.expectedExact, index.Exact)
		}
		if !reflect.DeepEqual(index.Glob, test.expectedGlob) {
			t.Fatalf("%s: expected glob index %+v, got %+v", test.description, test.expectedGlob, index.Glob)
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		lslintConfig *Config
		ignoreIndex  *IgnoreIndex
		path         string
		expected     bool
	}{
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: &IgnoreIndex{
				Exact: map[string]bool{
					".git": true,
				},
			},
			path:     ".git",
			expected: true,
		},
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: &IgnoreIndex{
				Exact: map[string]bool{
					"src": true,
				},
			},
			path:     "src/test/test.js",
			expected: true,
		},
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: &IgnoreIndex{
				Exact: map[string]bool{},
				Glob:  []string{".env*"},
			},
			path:     ".env.local",
			expected: true,
		},
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: &IgnoreIndex{
				Exact: map[string]bool{},
				Glob:  []string{"**/.env*"},
			},
			path:     "packages/ui/.env.local",
			expected: true,
		},
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: &IgnoreIndex{
				Exact: map[string]bool{},
				Glob:  []string{"packages/*/dist"},
			},
			path:     "packages/ui/dist/index.js",
			expected: true,
		},
	}

	i := 0
	for _, test := range tests {
		res := test.lslintConfig.ShouldIgnore(test.ignoreIndex, test.path)

		if res != test.expected {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}
