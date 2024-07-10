package main

import (
	"os"
	"strings"
)

const (
	NoPattern         = -1
	IfStatmentPattern = 0
	ForeachPattern    = 2
)

type State struct {
	Pattern int
	Info    map[string]any
}

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
		currentIndex           = -1
		ifStatmentsOpened      = 0
		forEachStatmentsOpened = 0
		currentVariableName    = ""
		content                = ""
		patternState           = State{}
		temlpateCopy           = strings.Clone(template)
	)

	for i := 0; i < len(template); i++ {
		switch template[i] {
		case '{':
			isInBracket = true
			currentVariableName = ""
			if ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
				content += "{"
				continue
			}

		case '}':
			isInBracket = false
			variableName := strings.Trim(currentVariableName, " ")
			claims := strings.Split(variableName, " ")

			if len(claims) == 1 {
				if variableName == "/if" {
					ifStatmentsOpened -= 1
					if ifStatmentsOpened != 0 {
						content += "}"
						continue
					}
 
					shouldRender := patternState.Info["shouldRender"].(bool)
					rest := template[i + 1:]
					
					if !shouldRender {
						temlpateCopy = temlpateCopy[0:currentIndex+1] + rest
						continue
					}

					content, _ := strings.CutSuffix(content, "{/if")
					content = t.RenderTemplate(content, variables)
					temlpateCopy = temlpateCopy[0:currentIndex+1] + content + rest
					content = ""
					currentIndex += len(content)
					patternState = State{}

				} else if variableName == "/foreach" {
					forEachStatmentsOpened -= 1
					continue
				}
				value := variables[variableName]
				if value != nil {
					strValue := Stringify(value)
					temlpateCopy = temlpateCopy[0:currentIndex+1] + strValue + template[i+1:]
					currentIndex += len(strValue)
				}
			} else {
				pattern := detectPattern(claims)
				if pattern == -1 || ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
					if pattern == IfStatmentPattern {
						ifStatmentsOpened += 1
					} else if pattern == ForeachPattern {
						forEachStatmentsOpened += 1
					}
					content += "}"
					continue
				}

				if pattern == IfStatmentPattern {
					ifStatmentsOpened += 1

					if ifStatmentsOpened != 1 {
						continue
					}

					value := variables[claims[1]]
					patternState = State{
						Pattern: IfStatmentPattern,
						Info: map[string]any{
							"shouldRender":  value.(bool),
							"ifStatmentLen": len(variableName) - 1,
						},
					}
				} else if pattern == ForeachPattern {
					forEachStatmentsOpened += 1

					if forEachStatmentsOpened != 1 {
						continue

					}
					patternState = State{
						Pattern: ForeachPattern,
						Info: map[string]any{
							"variableName":          claims[1],
							"iterationVariableName": claims[3],
						},
					}
				}
			}

		default:
			incrementCurrentIndex := true
			s := string(template[i])

			if ifStatmentsOpened > 0 || forEachStatmentsOpened > 0 {
				content += s
				incrementCurrentIndex = false
			}

			if isInBracket {
				currentVariableName += s
				incrementCurrentIndex = false
			}

			if incrementCurrentIndex {
				currentIndex++
			}
		}
	}

	return temlpateCopy
}

func detectPattern(claims []string) int {
	if isIfStatmentPattern(claims) {
		return IfStatmentPattern
	} else if isForeachPattern(claims) {
		return ForeachPattern
	}
	return -1
}

func isIfStatmentPattern(claims []string) bool {
	if len(claims) != 2 {
		return false
	}

	if claims[0] != "@if" {
		return false
	}

	return true
}

func isForeachPattern(claims []string) bool {
	if len(claims) != 4 {
		return false
	}

	if claims[0] != "@foreach" || claims[2] != "as" {
		return false
	}

	return true
}

//
// {@if isAdmin}
//     <h1>Jestem adminem</h1>
// {/if}
//
// {@foreach numbers as num}
//     <h1>{ num }</h1>
// {/foreach}
//
// {@if isAdmin}
// 	   { username }
//     <h1>Jestem adminem</h1>
// 	->
// 	   username
// 	   <h1>Jestem adminem</h1>
//
//
//
// {@if isAdmin}
//     <h1>Jestem adminem</h1>
//     {@if moreInfo}
// 		  <p>More Info</p>
//     {/if}
// 	   {@if isAdmin}
// 	       <h1>Jestem adminem</h1>
// 	       {@if moreInfo}
// 	   		  <p>More Info</p>
// 	       {/if}
// 	   {/if}
// {/if}
//
//
//
//
// if isAdmin {
//     t += render("<h1>Jestem adminem</h1>")
// }
//
