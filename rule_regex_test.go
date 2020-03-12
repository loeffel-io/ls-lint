package main

import "testing"

func TestRuleRegex(t *testing.T) {
	var tests = []*struct {
		params   []string
		value    string
		expected bool
		err      error
	}{
		{params: []string{".+"}, value: "regex", expected: true, err: nil},
		{params: []string{"[0-9]+"}, value: "123", expected: true, err: nil},
		{params: []string{"[a-z]+"}, value: "123", expected: false, err: nil},
		{params: []string{"[a-z\\-]+"}, value: "google-test", expected: true, err: nil},
		{params: []string{"[a-z\\-]+"}, value: "google.test", expected: false, err: nil},
	}

	var i = 0
	for _, test := range tests {
		var rule = new(RuleRegex).Init()

		// parameters
		err := rule.SetParameters(test.params)

		if err != nil && err != test.err {
			t.Errorf("Test %d failed with unmatched error - %s", i, err.Error())
		}

		// validate
		res, err := rule.Validate(test.value)

		if err != nil && err != test.err {
			t.Errorf("Test %d failed with unmatched error - %s", i, err.Error())
		}

		if res != test.expected {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
		}

		i++
	}
}
