package templateEngine

const (
	noStatmentCode = iota
	ifStatmentCode
	foreachStatmentCode
)

type StatmentState struct {
	Statment int
	Info     map[string]any
}

func getNewEmptyStatment() StatmentState {
	return StatmentState{
		Statment: noStatmentCode,
		Info:     make(map[string]any),
	}
}

func getNewStatment(statmentCode int, info map[string]any) StatmentState {
	return StatmentState{
		Statment: statmentCode,
		Info:     info,
	}
}
