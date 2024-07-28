package templateEngine

type statmentManagerInterface interface {

	GetStatmentCode() int

	GetStatmentEndingTag() string

	SetDependencies(statmentState *statmentState, variables *map[string]any)

	IsStatmentPatternCorrect(claims []string) bool

	RenderStatment(str string, renderTemplate renderFunction) string

	SetNewStatmentState(claims []string, statmentState *statmentState, claimsInBracket string)
}

type renderFunction func(template string, variables map[string]any) string
