package rule

import "testing"

func TestDisallow(t *testing.T) {
	var rule = new(Disallow).Init()

	var tests = []*ruleTest{
		{value: "camelCase", expected: false, err: nil},
		{value: "PascalCase", expected: false, err: nil},
		{value: "kebab-case", expected: false, err: nil},
		{value: "lowercase", expected: false, err: nil},
		{value: "point.case", expected: false, err: nil},
		{value: "snake_case", expected: false, err: nil},
		{value: "SCREAMING_SNAKE_CASE", expected: false, err: nil},
		{value: "literally anything at all", expected: false, err: nil},
	}

	var i = 0
	for _, test := range tests {
		res, err := rule.Validate(test.value)

		if err != nil && err != test.err {
			t.Errorf("Test %d failed with unmatched error - %s", i, err.Error())
			return
		}

		if res != test.expected {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, res)
			return
		}

		i++
	}
}
