package display

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"gowordle.com/display/handlers"
	"gowordle.com/display/model"
)

type State struct {
	State model.State
}

func NewModel() State {
	return State{State: model.State{Selected: model.ScreenHome}}
}

// State Init
func (m State) Init() tea.Cmd {
	return nil
}

// Update (auto-detect events)
func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.State.Editing {
			nm, cmd := handlers.HandleEditKey(m.State, msg)
			return State{State: nm}, cmd
		}
		nm, cmd := handlers.HandleKey(m.State, msg)
		return State{State: nm}, cmd
	case tea.WindowSizeMsg:
		nm, cmd := handlers.HandleWindowSize(m.State, msg)
		return State{State: nm}, cmd
	}
	return m, nil
}

// View (render UI)
func (m State) View() tea.View {
	theme := DefaultTheme()
	menu := m.menuView(theme)
	content := m.contentView(theme)
	help := theme.Muted.Render("Navigation: ↑ ↓ ← → / Tab, chiffres 1-5, Enter, q pour quitter")

	left := theme.Panel.Width(24).Render(menu)
	rightWidth := 66
	if m.State.Width > 30 {
		rightWidth = max(36, m.State.Width-32)
	}
	right := theme.Border.Width(rightWidth).Render(content)

	var bottom string
	if m.State.Height > 0 {
		bottom = help
	} else {
		bottom = help
	}

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right) + "\n\n" + bottom)
}

func (m State) menuView(theme Theme) string {
	items := make([]string, 0, len(model.ScreenLabels))
	for i, label := range model.ScreenLabels {
		style := theme.Muted
		prefix := "  "
		if model.Screen(i) == m.State.Selected {
			style = theme.Accent
			prefix = "> "
		}
		items = append(items, style.Render(prefix+fmt.Sprintf("%d. %s", i+1, label)))
	}

	return strings.Join(append([]string{theme.Section.Render("Menu")}, items...), "\n\n")
}

func (m State) contentView(theme Theme) string {
	switch m.State.Selected {
	case model.ScreenHome:
		return HomeScreen()
	case model.ScreenAuth:
		return AuthScreen(m.State)
	case model.ScreenSettings:
		return SettingsScreen()
	default:
		return HomeScreen()
	}
}

// Init tea program
func Run() error {
	program := tea.NewProgram(NewModel())
	_, err := program.Run()
	return err
}
