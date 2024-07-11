package rule

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestExists_GetParameters(t *testing.T) {
	var tests = []*struct {
		params   []string
		expected []string
		err      error
	}{
		{params: []string{"3"}, expected: []string{"3"}, err: nil},
		{params: []string{"0"}, expected: []string{"0"}, err: nil},
		{params: []string{"1-4"}, expected: []string{"1-4"}, err: nil},
		{params: []string{"-1"}, expected: []string{"0"}, err: strconv.ErrSyntax},
		{params: []string{"2342323423234"}, expected: []string{"0"}, err: strconv.ErrRange},
		{params: []string{"1-"}, expected: []string{"0"}, err: strconv.ErrSyntax},
		{params: []string{"1-2342323423234"}, expected: []string{"0"}, err: strconv.ErrRange},
	}

	var i = 0
	for _, test := range tests {
		var rule = new(Exists).Init()

		err := rule.SetParameters(test.params)
		if !errors.Is(err, test.err) {
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

func TestExists_Validate(t *testing.T) {
	var tests = []*struct {
		rule  *Exists
		fail  bool
		count uint16
		valid bool
		err   error
	}{
		{rule: &Exists{name: "exists", exclusive: true, min: 1, max: 1, count: 1, RWMutex: new(sync.RWMutex)}, fail: true, count: 1, valid: true, err: nil},
		{rule: &Exists{name: "exists", exclusive: true, min: 1, max: 3, count: 0, RWMutex: new(sync.RWMutex)}, fail: true, count: 0, valid: false, err: nil},
		{rule: &Exists{name: "exists", exclusive: true, min: 3, max: 6, count: 8, RWMutex: new(sync.RWMutex)}, fail: true, count: 8, valid: false, err: nil},
		{rule: &Exists{name: "exists", exclusive: true, min: 3, max: 6, count: 6, RWMutex: new(sync.RWMutex)}, fail: false, count: 7, valid: true, err: nil},
	}

	var i = 0
	for _, test := range tests {
		var valid, err = test.rule.Validate("", test.fail)

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
