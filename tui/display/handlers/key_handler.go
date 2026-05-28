package handlers

import (
	tea "charm.land/bubbletea/v2"
	"gowordle.com/display/handlers/fields"
	"gowordle.com/display/model"
)

func HandleKey(m model.State, msg tea.KeyMsg) (model.State, tea.Cmd) {
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
		if m.Selected == model.ScreenAuth {
			m.Auth.Field = model.AuthFieldLogin
		}
	}
	return m, nil
}

func HandleEditKey(m model.State, msg tea.KeyMsg) (model.State, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return exitEdit(m), nil
	case "enter":
		if m.Selected == model.ScreenAuth && m.Auth.Field == model.AuthFieldLogin {
			m.Auth.Field = model.AuthFieldPassword
			return m, nil
		}
		return exitEdit(m), nil
	case "tab", "down":
		if m.Selected == model.ScreenAuth {
			m.Auth.Field = fields.NextAuthField(m.Auth.Field)
			return m, nil
		}
		return exitEdit(m), nil
	case "backspace", "delete":
		if m.Selected == model.ScreenAuth {
			return deleteAuthValue(m), nil
		}
		return m, nil
	}

	if m.Selected == model.ScreenAuth && len([]rune(msg.String())) == 1 {
		return fields.AppendAuthValue(m, msg.String()), nil
	}
	return m, nil
}

func exitEdit(m model.State) model.State {
	m.Editing = false
	m.Focus = false
	return m
}

func deleteAuthValue(m model.State) model.State {
	return fields.DeleteAuthValue(m)
}

func prevScreen(m model.State) model.State {
	if m.Selected == 0 {
		m.Selected = model.Screen(len(model.ScreenLabels) - 1)
	} else {
		m.Selected--
	}
	return m
}

func nextScreen(m model.State) model.State {
	m.Selected = (m.Selected + 1) % model.Screen(len(model.ScreenLabels))
	return m
}

func selectByNumber(m model.State, s string) model.State {
	n := s[0] - '1'
	if n >= 0 && int(n) < len(model.ScreenLabels) {
		m.Selected = model.Screen(n)
	}
	return m
}
