package main

import "testing"

func TestReverseString(t *testing.T) {
	var statusesTest = []struct {
		input  string
		output string
	}{
		{"", ""},
		{"A", "A"},
		{"abcd", "dcba"},
		{"1234567890", "0987654321"},
	}

	for _, testCase := range statusesTest {
		expectedResult := testCase.output
		actualResult := reverse(testCase.input)
		if actualResult != testCase.output {
			t.Errorf("Failed UpdatePetStatus, expected %v, want %v", actualResult, expectedResult)
		}
	}

}
