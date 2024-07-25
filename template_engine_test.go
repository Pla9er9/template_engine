package main

import (
	"fmt"
	"testing"
)

var templateEngine = GetTemplateEngine()

func TestRenderVariable(t *testing.T) {
	var (
		variables = map[string]any{
			"username": "Alex",
			"color":    "blue",
			"number":   2,
			"array":    []string{"item1", "item2"},
			"boolean":  true,
		}
	)

	testCases := []TestCase{
		{
			Input:          "<h1>{username}{ username }{    username   }</h1>",
			ExpectedOutput: "<h1>AlexAlexAlex</h1>",
		},
		{
			Input:          "<h1 style='color: { color };'>{  username }</h1>",
			ExpectedOutput: "<h1 style='color: blue;'>Alex</h1>",
		},
		{
			Input:          "{ color }{  username }",
			ExpectedOutput: "blueAlex",
		},
		{
			Input:          "{ number } {  array } { boolean }",
			ExpectedOutput: "2 [\"item1\",\"item2\"] true",
		},
		{
			Input:          "{ empty }",
			ExpectedOutput: "{ empty }",
		},
		{
			Input:          "{ ello",
			ExpectedOutput: "{ ello",
		},
	}

	testRenderTestCases(t, testCases, variables)
}

func TestRenderIfStatment(t *testing.T) {
	variables := map[string]any{
		"isAuthorized": true,
		"isAdmin":      false,
		"moreInfo":     true,
		"message":      "Hello abc",
	}

	testCases := []TestCase{
		{
			Input:          "{@if isAuthorized}<h1>You are authorized</h1>{/if}",
			ExpectedOutput: "<h1>You are authorized</h1>",
		},
		{
			Input:          "{@if isAdmin}<h1>You are admin</h1>{/if}",
			ExpectedOutput: "",
		},
		{
			Input:          "{@if isAuthorized}<h1>You are authorized</h1>{@if moreInfo}<p>More Info</p>{/if}{@if moreInfo}<h1>Jestem adminem</h1>{@if isAdmin}<p>More Info</p>{/if}{/if}{/if}",
			ExpectedOutput: "<h1>You are authorized</h1><p>More Info</p><h1>Jestem adminem</h1>",
		},
		{
			Input:          "{@if isAuthorized}<h1>You are authorized</h1>{@if isAdmin}<p>More Info</p>{/if}{/if}",
			ExpectedOutput: "<h1>You are authorized</h1>",
		},
		{
			Input:          "{@if isAuthorized}{message}{/if}",
			ExpectedOutput: "Hello abc",
		},
		{
			Input:          "{@if none}{message}{/if}",
			ExpectedOutput: "{@if none}{message}{/if}",
		},
	}

	testRenderTestCases(t, testCases, variables)
}

func TestRenderForeachStatment(t *testing.T) {
	variables := map[string]any{
		"numbers": []int{1, 2, 3},
	}

	testCases := []TestCase{
		{
			Input:          "{@foreach numbers as num}{num}{/foreach}",
			ExpectedOutput: "123",
		},
		{
			Input:          "{@foreach numbers as num}lol{/foreach}",
			ExpectedOutput: "lollollol",
		},
		{
			Input:          "{@foreach numbers as num}{@foreach numbers as num}{num}{/foreach}{/foreach}",
			ExpectedOutput: "123123123",
		},
		{
			Input:          "{num{num}{/foreach",
			ExpectedOutput: "{num{num}{/foreach",
		},
	}

	testRenderTestCases(t, testCases, variables)
}

func testRenderTestCases(t *testing.T, testCases []TestCase, variables map[string]any) {
	for _, test := range testCases {
		input := test.Input.(string)
		output := templateEngine.RenderTemplate(input, variables)

		if output != test.ExpectedOutput {
			errorMsg := getRenderErrorMsg(test.ExpectedOutput.(string), output)
			t.Error(errorMsg)
		}
	}
}

func getRenderErrorMsg(expectedOutput string, output string) string {
	return fmt.Sprintf(
		"Expected output and final output are not same\nExpected output: \n%v \nFinal output: \n%v",
		expectedOutput,
		output,
	)
}
