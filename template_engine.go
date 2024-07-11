package main

import (
	"os"
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
		currentStatmentState   = StatmentState{}
		claimsInBracket        = ""
		contentInStatment      = ""
		result                 = ""
	)

	for i := 0; i < len(template); i++ {
		switch template[i] {
		case '{':
			isInBracket = true
			result += claimsInBracket
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

					shouldRender := currentStatmentState.Info["shouldRender"].(bool)
					if !shouldRender {
						continue
					}

					renderedContent := t.RenderTemplate(contentInStatment, variables)
					result += renderedContent
					contentInStatment = ""
					currentStatmentState = StatmentState{}

				} else if variableName == "/foreach" {
					forEachStatmentsOpened -= 1
				} else {
					value := variables[variableName]
					if value != nil {
						strValue := Stringify(value)
						result += strValue
					} else {
						result += "{" + claimsInBracketClone + "}"
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
					if ifStatmentsOpened == 1 {
						value := variables[claims[1]]
						currentStatmentState = GetNewIfStatmentState(value.(bool))
					}
				} else if statment == ForeachStatment {
					forEachStatmentsOpened++
					if forEachStatmentsOpened == 1 {
						currentStatmentState = GetNewForeachStatmentState(claims[1], claims[3])
					}
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
