package main

import "testing"

func TestRuleKebabCase(t *testing.T) {
	var rule = new(RuleKebabCase).Init()

	var tests = []*test{
		{value: "kebab", expected: true, err: nil},
		{value: "kebabcase", expected: true, err: nil},
		{value: "kebabCase", expected: false, err: nil},
		{value: "Kebabcase", expected: false, err: nil},
		{value: "KebabCase", expected: false, err: nil},
		{value: "KEBABCASE", expected: false, err: nil},
		{value: "kebab-case", expected: true, err: nil},
		{value: "kebab-case-test", expected: true, err: nil},
	}

	var i = 0
	for _, test := range tests {
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
