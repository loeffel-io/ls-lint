package rule

import (
	"errors"
	"testing"
)

func TestSnakeCase(t *testing.T) {
	rule := new(SnakeCase).Init()

	tests := []*ruleTest{
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
