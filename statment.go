package main

const (
	NoStatment      = -1
	IfStatment      = 0
	ForeachStatment = 2
)

type Statment struct {
	Statment int
	Info     map[string]any
}

func GetNewEmptyStatment() Statment {
	return Statment{
		Statment: NoStatment,
		Info: make(map[string]any),
	}
}

func GetNewIfStatment(shouldRender bool) Statment {
	return Statment{
		Statment: IfStatment,
		Info: map[string]any{
			"shouldRender": shouldRender,
		},
	}
}

func GetNewForeachStatment(variableName, iterationVariableName string) Statment {
	return  Statment{
		Statment: ForeachStatment,
		Info: map[string]any{
			"variableName":          variableName,
			"iterationVariableName": iterationVariableName,
		},
	}
}