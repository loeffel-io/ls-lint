package rule

import (
	"errors"
	"testing"
)

func TestPascalCase(t *testing.T) {
	rule := new(PascalCase).Init()

	tests := []*ruleTest{
		{value: "pascal", expected: false, err: nil},
		{value: "pascalcase", expected: false, err: nil},
		{value: "pascalCase", expected: false, err: nil},
		{value: "Pascalcase", expected: true, err: nil},
		{value: "PascalCase", expected: true, err: nil},
		{value: "PascałCase", expected: true, err: nil}, // here "l" has a diacritic mark: ł
		{value: "Pascal1Case", expected: true, err: nil},
		{value: "PascalVCase", expected: true, err: nil},
		{value: "PascalCaseForever", expected: true, err: nil},
		{value: "PASCALCASE", expected: false, err: nil},
		{value: "pascal_case", expected: false, err: nil},
		{value: "pascal.case", expected: false, err: nil},
		{value: "pascal-case", expected: false, err: nil},
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
