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
	help := theme.Muted.Render("Navigation: ↑ ↓ ← → / Tab, chiffres 1-4, Enter, q pour quitter")

	if m.State.Selected == model.ScreenGame {
		content := GameScreen(m.State)
		return tea.NewView(content + "\n\n" + help)
	}

	menu := m.menuView(theme)
	content := m.contentView()
	left := theme.Panel.Width(24).Render(menu)
	rightWidth := 66
	if m.State.Width > 30 {
		rightWidth = max(36, m.State.Width-32)
	}
	right := theme.Border.Width(rightWidth).Render(content)

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right) + "\n\n" + help)
}

func (m State) menuView(theme Theme) string {
	visible := visibleScreens(m.State)
	items := make([]string, 0, len(visible))
	for i, screen := range visible {
		style := theme.Muted
		prefix := "  "
		if screen == m.State.Selected {
			style = theme.Accent
			prefix = "> "
		}
		items = append(items, style.Render(prefix+fmt.Sprintf("%d. %s", i+1, model.ScreenLabels[screen])))
	}

	return strings.Join(append([]string{theme.Section.Render("Menu")}, items...), "\n\n")
}

func (m State) contentView() string {
	switch m.State.Selected {
	case model.ScreenHome:
		return HomeScreen()
	case model.ScreenAuth:
		return AuthScreen(m.State)
	case model.ScreenGame:
		return GameScreen(m.State)
	case model.ScreenSettings:
		return SettingsScreen()
	default:
		return HomeScreen()
	}
}

func visibleScreens(state model.State) []model.Screen {
	screens := []model.Screen{model.ScreenHome, model.ScreenAuth}
	if state.Connected {
		screens = append(screens, model.ScreenGame)
	}
	screens = append(screens, model.ScreenSettings)
	return screens
}

// Init tea program
func Run() error {
	program := tea.NewProgram(NewModel())
	_, err := program.Run()
	return err
}
