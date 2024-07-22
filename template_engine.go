package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type TemplateEngine struct{}

func getTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

func (t *TemplateEngine) RenderTemplateFromFile(filename string, variables map[string]any) (string, error) {
	str, err := readFile(filename)
	if err != nil {
		return "", err
	}
	return t.RenderTemplate(string(str), variables), nil
}

func readFile(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	return bytes, err
}

func (t *TemplateEngine) RenderTemplate(template string, variables map[string]any) string {
	var (
		isInBracket            = false
		ifStatmentsOpened      = 0
		forEachStatmentsOpened = 0
		currentStatmentState   = GetNewEmptyStatment()
		claimsInBracket   = ""
		contentInStatment = ""
		result            = ""
	)

	for i := 0; i < len(template); i++ {
		switch template[i] {
		case '{':
			if isInBracket {
				result += "{" + claimsInBracket
			}
			isInBracket = true
			claimsInBracket = ""
			if ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
				contentInStatment += "{"
				continue
			}

		case '}':
			isInBracket = false
			variableName := strings.Trim(claimsInBracket, " ")
			claims := strings.Split(variableName, " ")
			claimsInBracketClone := strings.Clone(claimsInBracket)
			claimsInBracket = ""

			if len(claims) == 1 {
				if variableName == "/if" {
					ifStatmentsOpened -= 1
					if ifStatmentsOpened != 0 {
						contentInStatment += "}"
						continue
					}

					shouldRender := currentStatmentState.Info["shouldRender"]

					if shouldRender == nil {
						ifStatment := currentStatmentState.Info["statment"].(string)
						result += "{" + ifStatment + "}" + contentInStatment + "}"
						continue
					}

					if !shouldRender.(bool) {
						continue
					}

					rep, ok := strings.CutSuffix(contentInStatment, "{/if")
					if ok {
						contentInStatment = rep
					}
					renderedContent := t.RenderTemplate(contentInStatment, variables)
					result += renderedContent
					contentInStatment = ""
					currentStatmentState = GetNewEmptyStatment()

				} else if variableName == "/foreach" {
					forEachStatmentsOpened -= 1
					if forEachStatmentsOpened != 0 {
						contentInStatment += "}"
						continue
					}

					variableName := currentStatmentState.Info["variableName"].(string)
					iterationVariableName := currentStatmentState.Info["iterationVariableName"].(string)
					array := variables[variableName]

					if array == nil {
						continue
					}

					slice, err := convertToSlice(array)

					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					for _, v := range slice {
						variables[iterationVariableName] = v
						contentInStatment, _ := strings.CutSuffix(contentInStatment, "{/foreach")
						result += t.RenderTemplate(contentInStatment, variables)
					}

					delete(variables, iterationVariableName)
					contentInStatment = ""
					currentStatmentState = GetNewEmptyStatment()

				} else {
					if ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
						contentInStatment += "}"
					} else {
						value := variables[variableName]
						if value != nil {
							strValue := Stringify(value)
							result += strValue
						} else {
							result += "{" + claimsInBracketClone + "}"
						}
					}
				}

			} else {
				statment := detectStatment(claims)
				if statment == -1 || ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
					if statment == IfStatment {
						ifStatmentsOpened += 1
					} else if statment == ForeachStatment {
						forEachStatmentsOpened += 1
					}
					contentInStatment += "}"
					continue
				}

				if statment == IfStatment {
					ifStatmentsOpened++
					if ifStatmentsOpened != 1 {
						continue
					}

					value := variables[claims[1]]
					valueType := reflect.TypeOf(value)

					if value == nil || valueType.String() != "bool" {
						currentStatmentState.Info["statment"] = claimsInBracketClone
						continue
					}
					currentStatmentState = GetNewIfStatment(value.(bool))

				} else if statment == ForeachStatment {
					forEachStatmentsOpened++
					if forEachStatmentsOpened != 1 {
						continue
					}

					value := variables[claims[1]]
					valueType := reflect.TypeOf(value)

					if value == nil || !strings.HasPrefix(valueType.String(), "[]") {
						currentStatmentState.Info["statment"] = claimsInBracketClone
						continue
					}

					currentStatmentState = GetNewForeachStatment(claims[1], claims[3])
				}

			}

		default:
			s := string(template[i])
			addToResult := true

			if ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
				contentInStatment += s
				addToResult = false
			}

			if isInBracket {
				claimsInBracket += s
				addToResult = false
			}

			if addToResult {
				result += s
			}
		}
	}

	if isInBracket {
		result += "{" + claimsInBracket
	}

	return result
}

func detectStatment(claims []string) int {
	if IsIfStatmentPattern(claims) {
		return IfStatment
	} else if IsForeachStatmentPattern(claims) {
		return ForeachStatment
	}
	return -1
}

func convertToSlice(data any) ([]any, error) {
	value := reflect.ValueOf(data)
	kind := value.Kind()

	if kind == reflect.Slice {
		sliceValue := value
		newSlice := make([]any, 0, sliceValue.Len())

		for i := 0; i < sliceValue.Len(); i++ {
			elementValue := sliceValue.Index(i)
			newSlice = append(newSlice, elementValue.Interface())
		}

		return newSlice, nil
	} else {
		return nil, fmt.Errorf("passed data is not slice")
	}
}
