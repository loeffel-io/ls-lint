package main

import "testing"

func TestRulePointCase(t *testing.T) {
	var rule = new(RulePointCase).Init()

	var tests = []*ruleTest{
		{value: "point", expected: true, err: nil},
		{value: "pointcase", expected: true, err: nil},
		{value: "pointCase", expected: false, err: nil},
		{value: "Pointcase", expected: false, err: nil},
		{value: "PointCase", expected: false, err: nil},
		{value: "POINTCASE", expected: false, err: nil},
		{value: "point.case", expected: true, err: nil},
		{value: "point.case.test", expected: true, err: nil},
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
