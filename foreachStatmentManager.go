package templateEngine

import (
	"fmt"
	"reflect"
	"strings"
)

func getForeachStatmentManager() *ForeachStatmentManager {
	return &ForeachStatmentManager{}
}

type ForeachStatmentManager struct {
	currentStatmentState *StatmentState
	variables            *map[string]any
}

func (s *ForeachStatmentManager) GetStatmentCode() int {
	return foreachStatmentCode
}

func (s *ForeachStatmentManager) GetStatmentEndingTag() string {
	return "/foreach"
}

func (s *ForeachStatmentManager) SetDependencies(statmentState *StatmentState, variables *map[string]any) {
	s.currentStatmentState = statmentState
	s.variables = variables
}

func (s *ForeachStatmentManager) IsStatmentPatternCorrect(claims []string) bool {
	if len(claims) != 4 {
		return false
	}

	if claims[0] != "@foreach" || claims[2] != "as" {
		return false
	}

	return true
}

func (s *ForeachStatmentManager) RenderStatment(str string, renderTemplate RenderFunction) string {
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

	array := (*s.variables)[variableName]
	if array == nil {
		return statmentStr.(string)
	}

	slice, err := convertAnyToSlice(array)

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

func (s *ForeachStatmentManager) SetNewStatmentState(claims []string, statmentState *StatmentState, claimsInBracketClone string) {
	var (
		value     = (*s.variables)[claims[1]]
		valueType = reflect.TypeOf(value)
	)

	if value == nil || !strings.HasPrefix(valueType.String(), "[]") {
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
