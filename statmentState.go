package templateEngine

const (
	noStatmentCode = iota
	ifStatmentCode
	foreachStatmentCode
)

type statmentState struct {
	Statment int
	Info     map[string]any
}

func getNewEmptyStatment() statmentState {
	return statmentState{
		Statment: noStatmentCode,
		Info:     make(map[string]any),
	}
}

func getNewStatment(statmentCode int, info map[string]any) statmentState {
	return statmentState{
		Statment: statmentCode,
		Info:     info,
	}
}
