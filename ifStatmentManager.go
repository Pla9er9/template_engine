package templateEngine

import (
	"reflect"
	"strings"
)

func getIfStatmentManager() *IfStatmentManager {
	return &IfStatmentManager{}
}

type IfStatmentManager struct {
	currentStatmentState *StatmentState
	variables            *map[string]any
}

func (s *IfStatmentManager) GetStatmentCode() int {
	return ifStatmentCode
}

func (s *IfStatmentManager) GetStatmentEndingTag() string {
	return "/if"
}

func (s *IfStatmentManager) SetDependencies(statmentState *StatmentState, variables *map[string]any) {
	s.currentStatmentState = statmentState
	s.variables = variables
}

func (s *IfStatmentManager) IsStatmentPatternCorrect(claims []string) bool {
	if len(claims) != 2 {
		return false
	}

	if claims[0] != "@if" {
		return false
	}

	return true
}

func (s *IfStatmentManager) RenderStatment(str string, renderTemplate RenderFunction) string {
	var (
		result       = ""
		shouldRender = s.currentStatmentState.Info["shouldRender"]
	)

	if shouldRender == nil {
		ifStatment := s.currentStatmentState.Info["statment"].(string)
		result += "{" + ifStatment + "}" + str + "}"
		return result
	}

	if !shouldRender.(bool) {
		return result
	}

	rep, ok := strings.CutSuffix(str, "{/if")
	if ok {
		str = rep
	}
	renderedContent := renderTemplate(str, (*s.variables))
	result += renderedContent

	return result
}

func (s *IfStatmentManager) SetNewStatmentState(claims []string, statmentState *StatmentState, claimsInBracketClone string) {
	var (
		value     = (*s.variables)[claims[1]]
		valueType = reflect.TypeOf(value)
	)

	if value == nil || valueType.String() != "bool" {
		statmentState.Info["statment"] = claimsInBracketClone
		return
	}
	
	*statmentState = getNewStatment(
		s.GetStatmentCode(),
		map[string]any{
			"shouldRender": value.(bool),
		},
	)
}
