package templateEngine

import "fmt"

func getMainStatmentManager(variables *map[string]any) *mainStatmentManager {
	var (
		emptyStatment          = getNewEmptyStatment()
		ifStatmentManager      = getIfStatmentManager()
		foreachStatmentManager = getForeachStatmentManager()
		managers               = []statmentManagerInterface{ifStatmentManager, foreachStatmentManager}
	)

	for _, m := range managers {
		m.SetDependencies(&emptyStatment, variables)
	}

	return &mainStatmentManager{
		statmentsOpened:       0,
		currentStatmentOpened: noStatmentCode,
		statmentState:         &emptyStatment,
		managers:              managers,
		variables:             variables,
	}
}

type mainStatmentManager struct {
	managers              []statmentManagerInterface
	statmentState         *statmentState
	statmentsOpened       int
	currentStatmentOpened int
	variables             *map[string]any
}

func (m *mainStatmentManager) ProcesStartingTag(claims []string) (matched bool) {
	for _, v := range m.managers {
		if !v.IsStatmentPatternCorrect(claims) {
			continue
		}

		if m.currentStatmentOpened != noStatmentCode && m.currentStatmentOpened != v.GetStatmentCode() {
			continue
		}

		m.statmentsOpened += 1
		m.currentStatmentOpened = v.GetStatmentCode()

		return true
	}
	return false
}

func (m *mainStatmentManager) ProcesEndingTag(endingTag string) (matched bool) {
	for _, v := range m.managers {
		if v.GetStatmentEndingTag() != endingTag {
			continue
		}

		if v.GetStatmentCode() == m.currentStatmentOpened || m.currentStatmentOpened == noStatmentCode {
			m.statmentsOpened -= 1
		}
		return true

	}

	return false
}

func (m *mainStatmentManager) SetNewStatmentState(claims []string, claimsInBracket string) {
	var currentManager statmentManagerInterface = *m.getStatmentManagerByCode(m.currentStatmentOpened)
	currentManager.SetNewStatmentState(claims, m.statmentState, claimsInBracket)
}

func (m *mainStatmentManager) ResetStatmentState() {
	*m.statmentState = getNewEmptyStatment()
	m.currentStatmentOpened = noStatmentCode
}

func (m *mainStatmentManager) RenderCurrentStatment(str string, renderTemplate renderFunction) string {
	var currentManager statmentManagerInterface = *m.getStatmentManagerByCode(m.currentStatmentOpened)
	return currentManager.RenderStatment(str, renderTemplate)
}

func (m *mainStatmentManager) getStatmentManagerByCode(statmentCode int) *statmentManagerInterface {
	if statmentCode > len(m.managers) || statmentCode == noStatmentCode {
		errMessage := fmt.Sprintf("Wrong statmentCode, statment code passed `%v`", statmentCode)
		panic(errMessage)
	}

	for _, v := range m.managers {
		if v.GetStatmentCode() == statmentCode {
			return &v
		}
	}

	return nil
}
