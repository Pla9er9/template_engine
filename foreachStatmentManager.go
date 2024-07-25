package main

import (
	"fmt"
	"reflect"
	"strings"
)

func GetForeachStatmentManager() *ForeachStatmentManager {
	return &ForeachStatmentManager{}
}

type ForeachStatmentManager struct {
	currentStatmentState *StatmentState
	variables            *map[string]any
}

func (s *ForeachStatmentManager) GetStatmentCode() int {
	return ForeachStatmentCode
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
		result                = ""
		variableName          = s.currentStatmentState.Info["variableName"].(string)
		iterationVariableName = s.currentStatmentState.Info["iterationVariableName"].(string)
		statmentStr           = fmt.Sprintf("{@foreach %v as %v}", variableName, iterationVariableName)
		array                 = (*s.variables)[variableName]
	)

	if array == nil {
		return statmentStr + str + "}"
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
		statmentState.Info["statment"] = claimsInBracketClone
		return
	}

	*statmentState = GetNewStatment(
		s.GetStatmentCode(),
		map[string]any{
			"variableName":          claims[1],
			"iterationVariableName": claims[3],
		},
	)
}
