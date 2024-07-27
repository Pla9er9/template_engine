package templateEngine

type StatmentManagerInterface interface {

	GetStatmentCode() int

	GetStatmentEndingTag() string

	SetDependencies(statmentState *StatmentState, variables *map[string]any)

	IsStatmentPatternCorrect(claims []string) bool

	RenderStatment(str string, renderTemplate RenderFunction) string

	SetNewStatmentState(claims []string, statmentState *StatmentState, claimsInBracket string)
}

type RenderFunction func(template string, variables map[string]any) string
