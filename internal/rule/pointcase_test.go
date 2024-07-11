package rule

import (
	"errors"
	"testing"
)

func TestPointCase(t *testing.T) {
	var rule = new(PointCase).Init()

	var tests = []*ruleTest{
		{value: "point", expected: true, err: nil},
		{value: "pointcase", expected: true, err: nil},
		{value: "pointCase", expected: false, err: nil},
		{value: "Pointcase", expected: false, err: nil},
		{value: "PointCase", expected: false, err: nil},
		{value: "POINTCASE", expected: false, err: nil},
		{value: "point.case", expected: true, err: nil},
		{value: "point12.case", expected: true, err: nil},
		{value: "point.case.test", expected: true, err: nil},
		{value: "point-case-test", expected: false, err: nil},
		{value: "point_case_test", expected: false, err: nil},
	}

	var i = 0
	for _, test := range tests {
		res, err := rule.Validate(test.value, true)

		if !errors.Is(err, test.err) {
			t.Errorf("Test %d failed with unmatched error - %e", i, err)
			return
		}

		if res != test.expected {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}
