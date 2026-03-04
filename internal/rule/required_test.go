package rule

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

func TestRequired_GetParameters(t *testing.T) {
	tests := []*struct {
		params   []string
		expected []string
		err      string
	}{
		{params: []string{"AGENTS.md"}, expected: []string{"AGENTS.md"}, err: ""},
		{params: []string{""}, expected: []string{""}, err: "required value is empty"},
	}

	i := 0
	for _, test := range tests {
		rule := new(Required).Init()

		err := rule.SetParameters(test.params)
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("Test %d failed with unmatched error - %e", i, err)
			return
		}

		params := rule.GetParameters()
		if !reflect.DeepEqual(params, test.expected) {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, params)
			return
		}

		i++
	}
}

func TestRequired_Validate(t *testing.T) {
	tests := []*struct {
		rule  *Required
		value string
		fail  bool
		count uint16
		valid bool
		err   error
	}{
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 0, RWMutex: new(sync.RWMutex)}, value: "README.md", fail: false, count: 0, valid: true, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 0, RWMutex: new(sync.RWMutex)}, value: "AGENTS.md", fail: false, count: 1, valid: true, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 0, RWMutex: new(sync.RWMutex)}, value: "", fail: true, count: 0, valid: false, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 1, RWMutex: new(sync.RWMutex)}, value: "", fail: true, count: 1, valid: true, err: nil},
	}

	i := 0
	for _, test := range tests {
		valid, err := test.rule.Validate(test.value, "", test.fail)

		if !errors.Is(err, test.err) {
			t.Errorf("Test %d failed with unmatched error - %e", i, err)
			return
		}

		if test.rule.count != test.count {
			t.Errorf("Test %d failed with unmatched count value - %d", i, test.rule.count)
			return
		}

		if valid != test.valid {
			t.Errorf("Test %d failed with unmatched return value - %+v", i, valid)
			return
		}

		i++
	}
}
