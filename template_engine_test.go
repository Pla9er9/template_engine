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
			Input:       "<h1>{username}{ username }{    username   }</h1>",
			ExpectedOutput: "<h1>" + username + username + username + "</h1>",
		},
		{
			Input:       "<h1 style='color: { color };'>{  username }</h1>",
			ExpectedOutput: "<h1 style='color: " + color + ";'>" + username + "</h1>",
		},
		{
			Input:       "{ color }{  username }",
			ExpectedOutput: color + username,
		},
		{
			Input:       "{ number } {  array } { boolean }",
			ExpectedOutput: "2 [\"item1\",\"item2\"] true",
		},
	}

	for _, test := range testCases {
		output, err := templateEngine.RenderTemplate(test.Input.(string), variables)

		if err != nil {
			t.Error(err.Error())
		}

		if output != test.ExpectedOutput {
			errorMsg := fmt.Sprintf(
				"Expected output and final output are not same\nExpected output: \n%v \nFinal output: \n%v",
				test.ExpectedOutput,
				output,
			)
			t.Error(errorMsg)
		}
	}
}
