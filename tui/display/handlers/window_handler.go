package handlers

import (
	tea "charm.land/bubbletea/v2"
	"gowordle.com/display/model"
)

func HandleWindowSize(m model.State, msg tea.WindowSizeMsg) (model.State, tea.Cmd) {
	m.Width = msg.Width
	m.Height = msg.Height
	return m, nil
}
