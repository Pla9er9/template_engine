package templateEngine

import (
	"reflect"
	"strings"
)

func getIfStatmentManager() *ifStatmentManager {
	return &ifStatmentManager{}
}

type ifStatmentManager struct {
	currentStatmentState *statmentState
	variables            *map[string]any
}

func (s *ifStatmentManager) GetStatmentCode() int {
	return ifStatmentCode
}

func (s *ifStatmentManager) GetStatmentEndingTag() string {
	return "/if"
}

func (s *ifStatmentManager) SetDependencies(statmentState *statmentState, variables *map[string]any) {
	s.currentStatmentState = statmentState
	s.variables = variables
}

func (s *ifStatmentManager) IsStatmentPatternCorrect(claims []string) bool {
	if len(claims) != 2 {
		return false
	}

	if claims[0] != "@if" {
		return false
	}

	return true
}

func (s *ifStatmentManager) RenderStatment(str string, renderTemplate renderFunction) string {
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

func (s *ifStatmentManager) SetNewStatmentState(claims []string, statmentState *statmentState, claimsInBracketClone string) {
	value, err := getVariable(claims[1], s.variables)
	valueType := reflect.TypeOf(value)

	if value == nil || valueType.String() != "bool" || err != nil {
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
