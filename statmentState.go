package main

const (
	NoStatmentCode = iota
	IfStatmentCode
	ForeachStatmentCode
)

type StatmentState struct {
	Statment int
	Info     map[string]any
}

func GetNewEmptyStatment() *StatmentState {
	return &StatmentState{
		Statment: NoStatmentCode,
		Info:     make(map[string]any),
	}
}

func GetNewStatment(statmentCode int, info map[string]any) StatmentState {
	return StatmentState{
		Statment: statmentCode,
		Info: info,
	}
}
