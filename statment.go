package main

const (
	NoStatment      = -1
	IfStatment      = 0
	ForeachStatment = 2
)

type StatmentState struct {
	Statment int
	Info     map[string]any
}

func GetNewIfStatmentState(shouldRender bool) StatmentState {
	return StatmentState{
		Statment: IfStatment,
		Info: map[string]any{
			"shouldRender": shouldRender,
		},
	}
}

func GetNewForeachStatmentState(variableName, iterationVariableName string) StatmentState {
	return  StatmentState{
		Statment: ForeachStatment,
		Info: map[string]any{
			"variableName":          variableName,
			"iterationVariableName": iterationVariableName,
		},
	}
}