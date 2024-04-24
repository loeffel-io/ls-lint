package rule

import (
	"testing"
)

func TestUpperCaseOrDigit(t *testing.T) {
	var rule = new(UppercaseDigit).Init()

	var tests = []*ruleTest{
		{value: "ALLCAPS", expected: true, err: nil},
		{value: "ALLCAPS007", expected: true, err: nil},
		{value: "VIN1234567", expected: true, err: nil},
		{value: "tooshortbud", expected: false, err: nil},
		{value: "SNEAKCase", expected: false, err: nil},
		{value: "Sneakcase", expected: false, err: nil},
		{value: "SneakCase", expected: false, err: nil},
		{value: "SNAKE_CASE", expected: false, err: nil},
		{value: "SNAKE_123_CASE", expected: false, err: nil},
		{value: "SNAKE_CASE_TEST", expected: false, err: nil},
		{value: "snake.case.test", expected: false, err: nil},
		{value: "SNAKE.CASE.TEST", expected: false, err: nil},
		{value: "snake-case-test", expected: false, err: nil},
		{value: "SNAKE-CASE-TEST", expected: false, err: nil},
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
