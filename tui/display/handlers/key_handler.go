package handlers

import (
	tea "charm.land/bubbletea/v2"
	"gowordle.com/display/handlers/fields"
	"gowordle.com/display/model"
)

func HandleKey(m model.State, msg tea.KeyMsg) (model.State, tea.Cmd) {
	if m.Selected == model.ScreenGame {
		return handleGameKey(m, msg)
	}

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

func handleGameKey(m model.State, msg tea.KeyMsg) (model.State, tea.Cmd) {
	if m.Game.Status != model.GamePlaying {
		if msg.String() == "esc" {
			m.Selected = model.ScreenHome
		}
		return m, nil
	}

	switch msg.String() {
	case "esc":
		m.Selected = model.ScreenHome
		return m, nil
	case "enter":
		return fields.SubmitGuess(m), nil
	case "backspace", "delete":
		return fields.DeleteGuessLetter(m), nil
	}

	if len([]rune(msg.String())) == 1 {
		r := []rune(msg.String())[0]
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return fields.AppendGuessLetter(m, msg.String()), nil
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
		if m.Selected == model.ScreenAuth && m.Auth.Field == model.AuthFieldPassword {
			m.Connected = true
			m.Selected = model.ScreenGame
			m.Game.WordToGuess = "BRUME"
			m.Editing = false
			m.Focus = false
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
	visible := visibleScreens(m)
	for i := range visible {
		if visible[i] == m.Selected {
			if i == 0 {
				m.Selected = visible[len(visible)-1]
			} else {
				m.Selected = visible[i-1]
			}
			return m
		}
	}
	m.Selected = visible[0]
	return m
}

func nextScreen(m model.State) model.State {
	visible := visibleScreens(m)
	for i := range visible {
		if visible[i] == m.Selected {
			m.Selected = visible[(i+1)%len(visible)]
			return m
		}
	}
	m.Selected = visible[0]
	return m
}

func selectByNumber(m model.State, s string) model.State {
	n := s[0] - '1'
	visible := visibleScreens(m)
	if n >= 0 && int(n) < len(visible) {
		m.Selected = visible[n]
	}
	return m
}

func visibleScreens(m model.State) []model.Screen {
	screens := []model.Screen{model.ScreenHome, model.ScreenAuth}
	if m.Connected {
		screens = append(screens, model.ScreenGame)
	}
	screens = append(screens, model.ScreenSettings)
	return screens
}
