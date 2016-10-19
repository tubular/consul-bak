package main_test

import (
	consulbak "github.com/Tubular/consul-bak"
	"testing"
)

type testCasePair struct {
	value  []string
	prefix string
	expect bool
}

var tests = []testCasePair{
	{[]string{"foo", "baz"}, "foobar", true},
	{[]string{"bar"}, "foobar", false},
}

func TestStartsWith(t *testing.T) {
	for _, testCase := range tests {
		result := consulbak.StartsWith(testCase.value, testCase.prefix)
		if result != testCase.expect {
			t.Errorf("For: %s in %s expected %t, actual: %t", testCase.value, testCase.prefix, testCase.expect, result)
		}
	}
}
