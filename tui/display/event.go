package display

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"gowordle.com/client"
	"gowordle.com/display/handlers"
	"gowordle.com/display/model"
)

type State struct {
	State  model.State
	Client *client.Client
}

func NewModel() State {
	cfg, _ := client.LoadConfig()
	c := client.New()

	s := model.State{Selected: model.ScreenHome}
	if cfg.Token != "" {
		c.SetToken(cfg.Token)
		s.Token = cfg.Token
		s.Connected = true
		s.Game.WordLoading = true
	}

	return State{State: s, Client: c}
}

func (m State) fetchWord() tea.Cmd {
	return func() tea.Msg {
		id, word, err := m.Client.RandomWord()
		return model.WordResultMsg{WordID: id, Word: word, Err: err}
	}
}

func (m State) saveGame() tea.Cmd {
	attempts, _ := json.Marshal(m.State.Game.TriedWords)
	duration := int(time.Since(m.State.Game.StartedAt).Seconds())
	wordID := m.State.Game.WordID
	won := m.State.Game.Status == model.GameWon
	return func() tea.Msg {
		err := m.Client.SaveGame(wordID, attempts, won, duration)
		return model.GameSaveResultMsg{Err: err}
	}
}

// State Init
func (m State) Init() tea.Cmd {
	if m.State.Connected {
		return m.fetchWord()
	}
	return nil
}

// Update (auto-detect events)
func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.State.Editing {
			nm, cmd := handlers.HandleEditKey(m.State, msg, func(email, password string) tea.Cmd {
				return func() tea.Msg {
					token, err := m.Client.Login(email, password)
					return model.AuthResultMsg{Token: token, Err: err}
				}
			})
			return State{State: nm, Client: m.Client}, cmd
		}
		nm, cmd := handlers.HandleKey(m.State, msg)
		var saveCmd tea.Cmd
		if m.State.Game.Status == model.GamePlaying && nm.Game.Status != model.GamePlaying {
			nm.Game.SaveLoading = true
			nm.Game.SaveError = ""
			saveCmd = State{State: nm, Client: m.Client}.saveGame()
		}
		return State{State: nm, Client: m.Client}, tea.Batch(cmd, saveCmd)
	case model.AuthResultMsg:
		nm := m.State
		nm.Auth.Loading = false
		if msg.Err != nil {
			nm.Auth.Error = msg.Err.Error()
			return State{State: nm, Client: m.Client}, nil
		}
		nm.Token = msg.Token
		nm.Connected = true
		nm.Selected = model.ScreenGame
		nm.Game.WordLoading = true
		nm.Auth.Error = ""
		nm.Auth.Login = ""
		nm.Auth.Password = ""
		m.Client.SetToken(msg.Token)
		_ = client.SaveConfig(&client.StoredConfig{Token: msg.Token})
		return State{State: nm, Client: m.Client}, m.fetchWord()
	case model.WordResultMsg:
		nm := m.State
		nm.Game.WordLoading = false
		if msg.Err == nil {
			nm.Game.WordID = msg.WordID
			nm.Game.WordToGuess = msg.Word
			nm.Game.StartedAt = time.Now()
		}
		return State{State: nm, Client: m.Client}, nil
	case model.GameSaveResultMsg:
		nm := m.State
		nm.Game.SaveLoading = false
		if msg.Err != nil {
			nm.Game.SaveError = msg.Err.Error()
		}
		return State{State: nm, Client: m.Client}, nil
	case tea.WindowSizeMsg:
		nm, cmd := handlers.HandleWindowSize(m.State, msg)
		return State{State: nm, Client: m.Client}, cmd
	}
	return m, nil
}

// View (render UI)
func (m State) View() tea.View {
	theme := DefaultTheme()
	help := theme.Muted.Render("Navigation: ↑ ↓ ← → / Tab, chiffres 1-4, Enter, q pour quitter")

	if m.State.Selected == model.ScreenGame {
		// Full screen
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
	screens := []model.Screen{model.ScreenHome}
	if state.Connected {
		screens = append(screens, model.ScreenGame)
	} else {
		screens = append(screens, model.ScreenAuth)
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
