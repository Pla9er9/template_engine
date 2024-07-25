package main

import (
	"os"
	"strings"
)

type TemplateEngine struct{}

func GetTemplateEngine() *TemplateEngine {
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
		statmentManager   = GetMainStatmentManager(&variables)
		isInBracket       = false
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

			if statmentManager.statmentsOpened > 0 {
				contentInStatment += "{"
				continue
			}

		case '}':
			var (
				symbol  = strings.Trim(claimsInBracket, " ")
				claims  = strings.Split(symbol, " ")
				cibCopy = strings.Clone(claimsInBracket)
			)

			claimsInBracket = ""
			isInBracket = false

			if len(claims) == 1 {
				matched := statmentManager.ProcesEndingTag(symbol)

				if statmentManager.statmentsOpened != 0 {
					contentInStatment += "}"
					continue
				}

				if matched {
					renderedStatment := statmentManager.RenderCurrentStatment(contentInStatment, t.RenderTemplate)
					result += renderedStatment
					statmentManager.ResetStatmentState()
					contentInStatment = ""

				} else {
					if statmentManager.statmentsOpened > 0 {
						contentInStatment += "}"
						continue
					}

					variableName := symbol
					value := variables[variableName]

					if value != nil {
						result += Stringify(value)
					} else {
						result += "{" + cibCopy + "}"
					}
				}

			} else {
				matched := statmentManager.ProcesStartingTag(claims)
				if statmentManager.statmentsOpened > 1 {
					contentInStatment += "}"
					continue
				}

				if !matched {
					result += "}"
					continue
				}

				statmentManager.SetNewStatmentState(claims, cibCopy)
			}

		default:
			s := string(template[i])
			addToResult := true

			if statmentManager.statmentsOpened > 0 {
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

	if statmentManager.statmentsOpened > 0 {
		result += contentInStatment
	} else if isInBracket {
		result += "{" + claimsInBracket
	}

	return result
}
