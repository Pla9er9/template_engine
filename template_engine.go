package templateEngine

import (
	"errors"
	"fmt"
	"os"
	"reflect"
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
		statmentManager   = getMainStatmentManager(&variables)
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

					result += t.renderVariable(cibCopy, &variables)
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

func (t *TemplateEngine) renderVariable(rawVariableName string, variables *map[string]any) string {
	var (
		value any
		err   error
		variableName     = strings.Trim(rawVariableName, " ")
		// example: user.name -> [user, name]
		properties       = strings.Split(variableName, ".")
		variableStatment = fmt.Sprintf("{%v}", rawVariableName)
	)

	if len(properties) > 1 {
		obj := (*variables)[properties[0]]
		value, err = t.getPropertyFromObject(obj, properties[1:])
	} else {
		value = (*variables)[variableName]
	}

	if value != nil && err == nil {
		return stringify(value)
	} else {
		return variableStatment
	}

}

func (t *TemplateEngine) getPropertyFromObject(struct_ any, fields []string) (any, error) {
	if struct_ == nil {
		return nil, errors.New("nil passed")
	}

	v, err := t.getField(&struct_, fields[0])
	if err != nil {
		return nil, err
	}

	if len(fields) == 1 {
		return v, nil
	}

	return t.getPropertyFromObject(v, fields[1:])
}

func (t *TemplateEngine) getField(v *any, field string) (any, error) {
	r := reflect.ValueOf(*v)
	f := reflect.Indirect(r).FieldByName(field)

	if f.Kind().String() == "invalid" || !f.CanInterface() {
		return nil, errors.New("field doesnt exist")
	}

	return f.Interface(), nil
}
