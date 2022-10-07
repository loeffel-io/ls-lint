package main

import (
	"bytes"
	"flag"
	"testing"
)

func TestMainRun(t *testing.T) {
	testSuite := []struct {
		description string
		args        []string
		expect      int
	}{
		{
			description: "Happy Path",
			args:        []string{"ls-lint"},
			expect:      0,
		},
	}

	for _, testCase := range testSuite {
		flag.Parse()
		t.Run(testCase.description, func(t *testing.T) {
			var output bytes.Buffer
			var actual = Run(&output, testCase.args)
			if actual != testCase.expect {
				t.Errorf("got %d but expected %d", actual, testCase.expect)
			}
		})
	}
}
