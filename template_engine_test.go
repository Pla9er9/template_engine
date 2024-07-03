package main

import (
	"fmt"
	"testing"
)

var templateEngine = getTemplateEngine()

func TestRenderVariable(t *testing.T) {
	username := "Alex"
	color := "blue"
	variables := map[string]any{
		"username": username,
		"color":    color,
		"number":   2,
		"array":    []string{"item1", "item2"},
		"boolean":  true,
	}

	testCases := []TestCase{
		{
			template:       "<h1>{username}{ username }{    username   }</h1>",
			expectedOutput: "<h1>" + username + username + username + "</h1>",
		},
		{
			template:       "<h1 style='color: { color };'>{  username }</h1>",
			expectedOutput: "<h1 style='color: " + color + ";'>" + username + "</h1>",
		},
		{
			template:       "{ color }{  username }",
			expectedOutput: color + username,
		},
		{
			template:       "{ number } {  array } { boolean }",
			expectedOutput: "2 [\"item1\",\"item2\"] true",
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
