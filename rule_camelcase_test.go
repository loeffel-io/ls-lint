package main

import "testing"

func TestRuleCamelCase(t *testing.T) {
	var rule = new(RuleCamelCase).Init()

	var tests = []*ruleTest{
		{value: "camel", expected: true, err: nil},
		{value: "camelcase", expected: true, err: nil},
		{value: "camelCase", expected: true, err: nil},
		{value: "camel1Case", expected: true, err: nil},
		{value: "camelVCase", expected: true, err: nil},
		{value: "camelCase123", expected: true, err: nil},
		{value: "camelCASE123", expected: false, err: nil},
		{value: "Camelcase", expected: false, err: nil},
		{value: "CamelCase", expected: false, err: nil},
		{value: "CAMELCASE", expected: false, err: nil},
		{value: "camel_case", expected: false, err: nil},
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
