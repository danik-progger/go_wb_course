package main

import (
	"testing"
)

type testCase struct {
	input    string
	expected string
	hasError bool
}

func TestStringUnpack(t *testing.T) {
	testCases := []testCase{
		// Standard cases from prompt
		{input: "a4bc2d5e", expected: "aaaabccddddde", hasError: false},
		{input: "abcd", expected: "abcd", hasError: false},

		// Escape sequences from prompt
		{input: `qwe\4\5`, expected: "qwe45", hasError: false},
		{input: `qwe\45`, expected: "qwe44444", hasError: false},
		{input: `qwe\\5`, expected: `qwe\\\\\`, hasError: false},

		// Invalid strings from prompt
		{input: "45", hasError: true},
		{input: "", expected: "", hasError: false},

		// Additional valid cases
		{input: "a10b", expected: "aaaaaaaaaab", hasError: false},
		{input: "a1b1c1", expected: "abc", hasError: false},
		{input: "a0b0c0", expected: "", hasError: false},
		{input: `\1\2\3`, expected: "123", hasError: false},
		{input: `\\`, expected: `\`, hasError: false},

		// Additional invalid cases
		{input: "3a", hasError: true},
		{input: `\`, hasError: true},
		{input: `abc\`, hasError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := StringUnpack(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("expected an error, but got none for input '%s'. result: %s", tc.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error for input '%s', but got: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("for input '%s', expected %q, but got %q", tc.input, tc.expected, result)
				}
			}
		})
	}
}
