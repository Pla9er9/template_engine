package main

import (
	"fmt"
	"testing"
)

var templateEngine = getTemplateEngine()

func TestRenderVariable(t *testing.T) {
	variables := map[string]any{
		"username": "Alex",
		"color":    "blue",
	}

	testCases := []TestCase{
		{
			template:       "<h1>{ username }</h1>",
			expectedOutput: "<h1>Alex</h1>",
		},
		{
			template:       "<h1>{    username   }</h1>",
			expectedOutput: "<h1>Alex</h1>",
		},
		{
			template:       "<h1 style='color: { color };'>{  username }</h1>",
			expectedOutput: "<h1 style='color: blue;'>Alex</h1>",
		},
		{
			template:       "{ color }{  username }",
			expectedOutput: "blueAlex",
		},
	}

	for _, test := range testCases {
		output, err := templateEngine.RenderTemplate(test.template, variables)

		if err != nil {
			t.Error(err.Error())
		}

		if output != test.expectedOutput {
			errorMsg := fmt.Sprintf(
				"Expected output and final output are not same\nExpected output: \n%v \nFinal output: \n%v",
				test.expectedOutput,
				output,
			)
			t.Error(errorMsg)
		}
	}
}
