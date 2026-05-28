package handlers

import (
	tea "charm.land/bubbletea/v2"
	"gowordle.com/display/state"
)

func HandleWindowSize(m state.Model, msg tea.WindowSizeMsg) (state.Model, tea.Cmd) {
	m.Width = msg.Width
	m.Height = msg.Height
	return m, nil
}
