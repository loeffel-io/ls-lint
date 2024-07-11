package config

import (
	"github.com/loeffel-io/ls-lint/v2/internal/rule"
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	var config = new(Config)
	var indexMock = map[string]map[string][]rule.Rule{
		".": {
			".dir": []rule.Rule{rule.RulesIndex["lowercase"]},
		},
		"./src": {
			".dir": []rule.Rule{rule.RulesIndex["camelcase"]},
		},
	}
	var indexMockEmpty = make(RuleIndex)

	var tests = []*struct {
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

	var i = 0
	for _, test := range tests {
		_, res := test.config.GetConfig(test.index, test.path)

		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		lslintConfig *Config
		ignoreIndex  map[string]bool
		path         string
		expected     bool
	}{
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: map[string]bool{
				".git": true,
			},
			path:     ".git",
			expected: true,
		},
		{
			lslintConfig: NewConfig(nil, nil),
			ignoreIndex: map[string]bool{
				"src": true,
			},
			path:     "src/test/test.js",
			expected: true,
		},
	}

	var i = 0
	for _, test := range tests {
		res := test.lslintConfig.ShouldIgnore(test.ignoreIndex, test.path)

		if res != test.expected {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}
