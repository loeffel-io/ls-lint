package rule

import "testing"

func TestSnakeCase(t *testing.T) {
	var rule = new(SnakeCase).Init()

	var tests = []*ruleTest{
		{value: "sneak", expected: true, err: nil},
		{value: "sneakcase", expected: true, err: nil},
		{value: "sneakCase", expected: false, err: nil},
		{value: "Sneakcase", expected: false, err: nil},
		{value: "SneakCase", expected: false, err: nil},
		{value: "SNEAKCASE", expected: false, err: nil},
		{value: "snake_case", expected: true, err: nil},
		{value: "snake_123_case", expected: true, err: nil},
		{value: "snake_case_test", expected: true, err: nil},
		{value: "snake.case.test", expected: false, err: nil},
		{value: "snake-case-test", expected: false, err: nil},
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
