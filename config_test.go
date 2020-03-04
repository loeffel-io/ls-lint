package main

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	var config = new(Config)
	var indexMock = map[string]map[string][]Rule{
		".": {
			".dir": []Rule{definitions["lowercase"]},
		},
		"./src": {
			".dir": []Rule{definitions["camelcase"]},
		},
	}
	var indexMockEmpty = map[string]map[string][]Rule{}

	var tests = []*struct {
		config   *Config
		index    index
		path     string
		expected map[string][]Rule
	}{
		{
			config: config,
			index:  indexMock,
			path:   "./src/test/Test.js",
			expected: map[string][]Rule{
				".dir": {definitions["camelcase"]},
			},
		},
		{
			config: config,
			index:  indexMock,
			path:   "./images/path.png",
			expected: map[string][]Rule{
				".dir": {definitions["lowercase"]},
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
		res := test.config.getConfig(test.index, test.path)

		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
		}

		i++
	}
}
