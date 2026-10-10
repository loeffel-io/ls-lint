package rule

import (
	"errors"
	"testing"
)

func TestCamelCase(t *testing.T) {
	rule := new(CamelCase).Init()

	tests := []*ruleTest{
		{value: "camel", expected: true, err: nil},
		{value: "camelcase", expected: true, err: nil},
		{value: "camelCase", expected: true, err: nil},
		{value: "camel1Case", expected: true, err: nil},
		{value: "camelVCase", expected: true, err: nil},
		{value: "camelCase123", expected: true, err: nil},
		{value: "camelCaseForever", expected: true, err: nil},
		{value: "camelCASE123", expected: false, err: nil},
		{value: "Camelcase", expected: false, err: nil},
		{value: "CamelCase", expected: false, err: nil},
		{value: "CAMELCASE", expected: false, err: nil},
		{value: "camel_case", expected: false, err: nil},
		{value: "camel.case", expected: false, err: nil},
	}

	i := 0
	for _, test := range tests {
		res, err := rule.Validate(test.value, "", true)

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
