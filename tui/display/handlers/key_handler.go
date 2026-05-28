package handlers

import (
	tea "charm.land/bubbletea/v2"
	"gowordle.com/display/state"
)

func HandleKey(m state.Model, msg tea.KeyMsg) (state.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "left", "up", "k":
		return prevScreen(m), nil
	case "right", "down", "j", "tab":
		return nextScreen(m), nil
	case "1", "2", "3":
		return selectByNumber(m, msg.String()), nil
	case "enter":
		m.Editing = true
		m.Focus = true
		if m.Selected == state.ScreenAuth {
			m.AuthField = state.AuthFieldLogin
		}
	}
	return m, nil
}

func HandleEditKey(m state.Model, msg tea.KeyMsg) (state.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return exitEdit(m), nil
	case "enter":
		if m.Selected == state.ScreenAuth && m.AuthField == state.AuthFieldLogin {
			m.AuthField = state.AuthFieldPassword
			return m, nil
		}
		return exitEdit(m), nil
	case "tab", "down":
		if m.Selected == state.ScreenAuth {
			m.AuthField = nextAuthField(m.AuthField)
			return m, nil
		}
		return exitEdit(m), nil
	case "backspace", "delete":
		if m.Selected == state.ScreenAuth {
			return deleteAuthValue(m), nil
		}
		return m, nil
	}

	if m.Selected == state.ScreenAuth && len([]rune(msg.String())) == 1 {
		return appendAuthValue(m, msg.String()), nil
	}
	return m, nil
}

func exitEdit(m state.Model) state.Model {
	m.Editing = false
	m.Focus = false
	return m
}

func nextAuthField(field state.AuthField) state.AuthField {
	if field == state.AuthFieldLogin {
		return state.AuthFieldPassword
	}
	return state.AuthFieldLogin
}

func appendAuthValue(m state.Model, value string) state.Model {
	if value == "" {
		return m
	}

	switch m.AuthField {
	case state.AuthFieldPassword:
		m.Password += value
	default:
		m.Login += value
	}

	return m
}

func deleteAuthValue(m state.Model) state.Model {
	switch m.AuthField {
	case state.AuthFieldPassword:
		if len(m.Password) > 0 {
			m.Password = m.Password[:len(m.Password)-1]
		}
	default:
		if len(m.Login) > 0 {
			m.Login = m.Login[:len(m.Login)-1]
		}
	}

	return m
}

func prevScreen(m state.Model) state.Model {
	if m.Selected == 0 {
		m.Selected = state.Screen(len(state.ScreenLabels) - 1)
	} else {
		m.Selected--
	}
	return m
}

func nextScreen(m state.Model) state.Model {
	m.Selected = (m.Selected + 1) % state.Screen(len(state.ScreenLabels))
	return m
}

func selectByNumber(m state.Model, s string) state.Model {
	n := s[0] - '1'
	if n >= 0 && int(n) < len(state.ScreenLabels) {
		m.Selected = state.Screen(n)
	}
	return m
}
