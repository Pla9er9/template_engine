package main

import (
	"encoding/json"
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
	return t.RenderTemplate(string(str), variables)
}

func readFile(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	return bytes, err
}

func (t *TemplateEngine) RenderTemplate(template string, variables map[string]any) (string, error) {
	var (
		isInBracket = false
		currentIndex = -1
		currentVariableName = ""
		temlpateCopy = strings.Clone(template)
	)

	for i := 0; i < len(template); i++ {
		switch template[i] {
		case '{':
			isInBracket = true
			currentVariableName = ""
			
		case '}':
			isInBracket = false
			variableName := strings.Trim(currentVariableName, " ")
			value := variables[variableName]
			
			if value != nil {
				var strValue string
				switch value.(type) {
				case string:
					strValue = value.(string)
				default:
					b, _ := json.Marshal(value)
					strValue = string(b)
				}
				temlpateCopy = temlpateCopy[0:currentIndex + 1] + strValue + template[i + 1:]
				currentIndex += len(strValue)
			}
			
		default:
			if (isInBracket) {
				currentVariableName += string(template[i])
			} else {
				currentIndex++
			}
		}
	}

	return temlpateCopy, nil
}
