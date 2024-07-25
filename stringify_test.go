package templateEngine

import (
	"fmt"
	"testing"
)

func TestStringify(t *testing.T) {
	testCases := []testCase{
		{
			Input:          "aba",
			ExpectedOutput: "aba",
		},
		{
			Input:          1,
			ExpectedOutput: "1",
		},
		{
			Input:          true,
			ExpectedOutput: "true",
		},
		{
			Input:          []string{"item1", "item2"},
			ExpectedOutput: "[\"item1\",\"item2\"]",
		},
		{
			Input: map[string]int{
				"water": 2,
				"oil":   3,
				"milk":  4,
			},
			ExpectedOutput: "{\"milk\":4,\"oil\":3,\"water\":2}",
		},
	}

	for _, test := range testCases {
		str := stringify(test.Input)

		if str != test.ExpectedOutput {
			errorMsg := fmt.Sprintf(
				"Stringify gave wrong output\nExpected output: \n%v \nFinal output: \n%v",
				test.ExpectedOutput,
				str,
			)
			t.Error(errorMsg)
		}
	}
}
