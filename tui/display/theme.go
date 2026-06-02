package display

import (
	"strings"

	"charm.land/lipgloss/v2"
)

type GameLetters struct {
	Correct   lipgloss.Style
	Misplaced lipgloss.Style
	Incorrect lipgloss.Style
}

// Theme
type Theme struct {
	Border  lipgloss.Style
	Title   lipgloss.Style
	Section lipgloss.Style
	Text    lipgloss.Style
	Accent  lipgloss.Style
	Muted   lipgloss.Style
	Error   lipgloss.Style
	Panel   lipgloss.Style
	Button  lipgloss.Style
	Letters GameLetters
}

func DefaultTheme() Theme {
	return Theme{
		Border:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(1, 2),
		Title:   lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true).Align(lipgloss.Center),
		Section: lipgloss.NewStyle().Foreground(lipgloss.Color("111")).Bold(true),
		Text:    lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Accent:  lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true),
		Muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Error:   lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
		Panel:   lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("236")).Padding(1, 2),
		Button:  lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Padding(0, 2).Bold(true),
		Letters: GameLetters{
			Correct:   lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("34")),
			Misplaced: lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("220")),
			Incorrect: lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("240")),
		},
	}
}

func RenderApp(title, subtitle string, body []string, footer string) string {
	theme := DefaultTheme()

	var rows []string
	rows = append(rows, theme.Title.Width(58).Render(title))
	if subtitle != "" {
		rows = append(rows, theme.Muted.Render(subtitle))
	}
	if len(body) > 0 {
		rows = append(rows, "")
		rows = append(rows, strings.Join(body, "\n"))
	}
	if footer != "" {
		rows = append(rows, "")
		rows = append(rows, theme.Muted.Render(footer))
	}

	content := strings.Join(rows, "\n")
	return theme.Border.Width(64).Render(content)
}
