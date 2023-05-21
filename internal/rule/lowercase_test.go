package rule

import "testing"

func TestLowercase(t *testing.T) {
	var rule = new(Lowercase).Init()

	var tests = []*ruleTest{
		{value: "abC", expected: false, err: nil},
		{value: "abc", expected: true, err: nil},
		{value: "abc-1", expected: true, err: nil},
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
