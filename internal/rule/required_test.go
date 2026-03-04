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
		{params: []string{}, expected: []string{}, err: ""},
		{params: []string{"AGENTS.md"}, expected: []string{"AGENTS.md"}, err: ""},
		{params: []string{""}, expected: []string{}, err: "required value is empty"},
	}

	testIndex := 0
	for _, test := range tests {
		rule := new(Required).Init()

		err := rule.SetParameters(test.params)
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("Test %d failed with unmatched error - %v", testIndex, err)
			return
		}

		params := rule.GetParameters()
		if !reflect.DeepEqual(params, test.expected) {
			t.Errorf("Test %d failed with unmatched return value - %+v", testIndex, params)
			return
		}

		testIndex++
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
		{rule: &Required{name: "required", exclusive: true, value: "", count: 0, RWMutex: new(sync.RWMutex)}, value: "README.md", fail: false, count: 1, valid: true, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 0, RWMutex: new(sync.RWMutex)}, value: "README.md", fail: false, count: 0, valid: true, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 0, RWMutex: new(sync.RWMutex)}, value: "AGENTS.md", fail: false, count: 1, valid: true, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "", count: 0, RWMutex: new(sync.RWMutex)}, value: "", fail: true, count: 0, valid: false, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 0, RWMutex: new(sync.RWMutex)}, value: "", fail: true, count: 0, valid: false, err: nil},
		{rule: &Required{name: "required", exclusive: true, value: "AGENTS.md", count: 1, RWMutex: new(sync.RWMutex)}, value: "", fail: true, count: 1, valid: true, err: nil},
	}

	testIndex := 0
	for _, test := range tests {
		valid, err := test.rule.Validate(test.value, "", test.fail)

		if !errors.Is(err, test.err) {
			t.Errorf("Test %d failed with unmatched error - %v", testIndex, err)
			return
		}

		if test.rule.count != test.count {
			t.Errorf("Test %d failed with unmatched count value - %d", testIndex, test.rule.count)
			return
		}

		if valid != test.valid {
			t.Errorf("Test %d failed with unmatched return value - %+v", testIndex, valid)
			return
		}

		testIndex++
	}
}

func TestRequired_SetParametersErrorDoesNotMutateValue(t *testing.T) {
	r := new(Required).Init()

	if err := r.SetParameters([]string{"AGENTS.md"}); err != nil {
		t.Fatalf("unexpected error setting initial parameter: %v", err)
	}

	if err := r.SetParameters([]string{""}); err == nil {
		t.Fatal("expected error for empty required value")
	}

	params := r.GetParameters()
	if !reflect.DeepEqual(params, []string{"AGENTS.md"}) {
		t.Fatalf("expected parameters to remain unchanged after error, got: %+v", params)
	}
}
