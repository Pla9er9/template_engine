package templateEngine

import (
	"fmt"
	"reflect"
	"strings"
)

func getForeachStatmentManager() *foreachStatmentManager {
	return &foreachStatmentManager{}
}

type foreachStatmentManager struct {
	currentStatmentState *statmentState
	variables            *map[string]any
}

func (s *foreachStatmentManager) GetStatmentCode() int {
	return foreachStatmentCode
}

func (s *foreachStatmentManager) GetStatmentEndingTag() string {
	return "/foreach"
}

func (s *foreachStatmentManager) SetDependencies(statmentState *statmentState, variables *map[string]any) {
	s.currentStatmentState = statmentState
	s.variables = variables
}

func (s *foreachStatmentManager) IsStatmentPatternCorrect(claims []string) bool {
	if len(claims) != 4 {
		return false
	}

	if claims[0] != "@foreach" || claims[2] != "as" {
		return false
	}

	return true
}

func (s *foreachStatmentManager) RenderStatment(str string, renderTemplate renderFunction) string {
	var (
		result                 = ""
		variableName_          = s.currentStatmentState.Info["variableName"]
		iterationVariableName_ = s.currentStatmentState.Info["iterationVariableName"]
		statmentStr            = s.currentStatmentState.Info["statment"]
	)

	if variableName_ == nil || iterationVariableName_ == nil {
		return statmentStr.(string) + str + "}"
	}

	variableName := variableName_.(string)
	iterationVariableName := iterationVariableName_.(string)

	value, err := getVariable(variableName, s.variables)
	if value == nil || err != nil {
		return statmentStr.(string)
	}

	slice, err := convertAnyToSlice(value)

	if err != nil {
		fmt.Printf("%v - %v\n", statmentStr, err.Error())
		return ""
	}

	val := (*s.variables)[iterationVariableName]

	for _, v := range slice {
		(*s.variables)[iterationVariableName] = v
		str, _ := strings.CutSuffix(str, "{/foreach")
		result += renderTemplate(str, (*s.variables))
	}

	if val == nil {
		delete((*s.variables), iterationVariableName)
	} else {
		(*s.variables)[iterationVariableName] = val
	}

	return result
}

func (s *foreachStatmentManager) SetNewStatmentState(claims []string, statmentState *statmentState, claimsInBracketClone string) {
	value, err := getVariable(claims[1], s.variables)
	valueType := reflect.TypeOf(value)

	if value == nil || !strings.HasPrefix(valueType.String(), "[]") || err != nil {
		statmentState.Info["statment"] = fmt.Sprintf("{@foreach %v as %v}", claims[1], claims[3])
		return
	}

	*statmentState = getNewStatment(
		s.GetStatmentCode(),
		map[string]any{
			"variableName":          claims[1],
			"iterationVariableName": claims[3],
		},
	)
}
