package templateEngine

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

	testCases := []testCase{
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
		{
			Input: `
			<style>
				button {
					min-width: 120px;
				}
			</style>`,
			ExpectedOutput: `
			<style>
				button {
					min-width: 120px;
				}
			</style>`,
		},
		{
			Input: `function deleteRow(id) {
            document.getElementById('id-${id}').remove()
        }`,
			ExpectedOutput: `function deleteRow(id) {
            document.getElementById('id-${id}').remove()
        }`,
		},
		{
			Input: `function deleteRow(id) {
            () => {document.getElementById('id-${id}').remove()}
        }`,
			ExpectedOutput: `function deleteRow(id) {
            () => {document.getElementById('id-${id}').remove()}
        }`,
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
		"array":        []string{"item1", "item2"},
	}

	testCases := []testCase{
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
			Input:          "{@if isAuthorized}{@foreach array as arr}{arr}{/foreach}{/if}",
			ExpectedOutput: "item1item2",
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
		"var":     true,
	}

	testCases := []testCase{
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
			Input:          "{@foreach numbers as num}{@if var}var{/if}{/foreach}",
			ExpectedOutput: "varvarvar",
		},
		{
			Input:          "{num{num}{/foreach",
			ExpectedOutput: "{num{num}{/foreach",
		},
	}

	testRenderTestCases(t, testCases, variables)
}

func TestRenderObjectVariables(t *testing.T) {
	type user struct {
		Name    string
		IsAdult bool
		Hobbies []string
		age     int
	}

	variables := map[string]any{
		"user": user{
			Name:    "Alex",
			age:     20,
			IsAdult: true,
			Hobbies: []string{
				"a", "b", "c",
			},
		},
		"lol": []string{
			"a", "b", "c",
		},
	}

	testCases := []testCase{
		{
			Input:          "{user.Name}",
			ExpectedOutput: "Alex",
		},
		// unexported field wont be rendered
		{
			Input:          "{user.age}",
			ExpectedOutput: "{user.age}",
		},
		{
			Input:          "{user.noexist}",
			ExpectedOutput: "{user.noexist}",
		},
		{
			Input:          "{@if user.IsAdult}Adult{/if}",
			ExpectedOutput: "Adult",
		},
		{
			Input:          "{@foreach user.Hobbies as h}{h}{/foreach}",
			ExpectedOutput: "abc",
		},
		{
			Input:          "{@foreach user.Hobbies as h}{h}{/foreach}",
			ExpectedOutput: "abc",
		},
	}

	testRenderTestCases(t, testCases, variables)
}

func testRenderTestCases(t *testing.T, testCases []testCase, variables map[string]any) {
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
