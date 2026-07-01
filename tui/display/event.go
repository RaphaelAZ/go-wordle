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
	"gowordle.com/display/lang"
	"gowordle.com/display/model"
)

type State struct {
	State  model.State
	Client *client.Client
}

func NewModel() State {
	cfg, _ := client.LoadConfig()
	c := client.New(cfg.Settings.BackendUrl)

	s := model.State{Selected: model.ScreenHome}
	s.Settings = storedToSettings(cfg.Settings)
	lang.Init(s.Settings.Language)
	SetTheme(s.Settings.Theme)
	if cfg.Token != "" {
		c.SetToken(cfg.Token)
		s.Token = cfg.Token
		s.Connected = true
		s.Game.WordLoading = true
	}

	return State{State: s, Client: c}
}

func storedToSettings(s client.StoredSettings) model.Settings {
	settings := model.DefaultSettings()
	if s.Theme != "" {
		settings.Theme = s.Theme
	}
	if s.Language != "" {
		settings.Language = s.Language
	}
	if s.DisplayMode != "" {
		settings.DisplayMode = s.DisplayMode
	}
	if s.BackendUrl != "" {
		settings.BackendUrl = s.BackendUrl
	}
	return settings
}

func (m State) currentStoredConfig() *client.StoredConfig {
	return &client.StoredConfig{
		Token: m.State.Token,
		Settings: client.StoredSettings{
			Theme:       m.State.Settings.Theme,
			Language:    m.State.Settings.Language,
			DisplayMode: m.State.Settings.DisplayMode,
			BackendUrl:  m.State.Settings.BackendUrl,
		},
	}
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
		updated := State{State: nm, Client: m.Client}
		_ = client.SaveConfig(updated.currentStoredConfig())
		return updated, updated.fetchWord()
	case model.WordResultMsg:
		nm := m.State
		nm.Game.WordLoading = false
		if msg.Err == nil {
			nm.Game.WordID = msg.WordID
			nm.Game.WordToGuess = msg.Word
			nm.Game.WordError = ""
			nm.Game.StartedAt = time.Now()
		} else {
			nm.Game.WordError = msg.Err.Error()
		}
		return State{State: nm, Client: m.Client}, nil
	case model.RestartGameMsg:
		return m, m.fetchWord()
	case model.LogoutMsg:
		nm := m.State
		nm.Token = ""
		nm.Connected = false
		nm.Game = model.Game{}
		nm.Selected = model.ScreenAuth
		nm.Settings.Field = model.SettingsFieldTheme
		m.Client.SetToken("")
		updated := State{State: nm, Client: m.Client}
		_ = client.SaveConfig(updated.currentStoredConfig())
		return updated, nil
	case model.SettingsChangedMsg:
		lang.Init(m.State.Settings.Language)
		SetTheme(m.State.Settings.Theme)
		_ = client.SaveConfig(m.currentStoredConfig())
		return m, nil
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
	help := theme.Muted.Render(lang.T("nav_help"))

	if m.State.Selected == model.ScreenGame {
		if m.State.Settings.DisplayMode == "normal" {
			nav := m.topNavView(theme)
			// Réduire la hauteur de 2 (nav + barre d'aide) pour éviter le débordement
			gs := m.State
			gs.Height = max(10, m.State.Height-2)
			content := GameScreen(gs)
			return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, nav, content, help))
		}
		// compact: side menu + game content (RenderApp already adds the border)
		content := GameScreen(m.State)
		menu := m.menuView(theme)
		left := theme.Panel.Width(24).Render(menu)
		return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", content) + "\n\n" + help)
	}

	if m.State.Settings.DisplayMode == "normal" {
		nav := m.topNavView(theme)
		content := m.contentView()
		w := m.State.Width
		if w <= 0 {
			w = 100
		}
		body := theme.Border.Width(w - 4).Render(content)
		return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, nav, body, help))
	}

	// compact: side panel + content
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
		items = append(items, style.Render(prefix+fmt.Sprintf("%d. %s", i+1, screenLabel(screen))))
	}

	return strings.Join(append([]string{theme.Section.Render(lang.T("menu_title"))}, items...), "\n\n")
}

func (m State) topNavView(theme Theme) string {
	visible := visibleScreens(m.State)
	items := make([]string, 0, len(visible))
	for i, screen := range visible {
		label := fmt.Sprintf(" %d. %s ", i+1, screenLabel(screen))
		if screen == m.State.Selected {
			items = append(items, theme.Button.Render(label))
		} else {
			items = append(items, theme.Muted.Render(label))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, items...)
}

func (m State) contentView() string {
	switch m.State.Selected {
	case model.ScreenHome:
		return HomeScreen(m.State)
	case model.ScreenAuth:
		return AuthScreen(m.State)
	case model.ScreenGame:
		return GameScreen(m.State)
	case model.ScreenSettings:
		return SettingsScreen(m.State)
	default:
		return HomeScreen(m.State)
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

func screenLabel(s model.Screen) string {
	switch s {
	case model.ScreenHome:
		return lang.T("screen_home")
	case model.ScreenAuth:
		return lang.T("screen_auth")
	case model.ScreenGame:
		return lang.T("screen_game")
	case model.ScreenSettings:
		return lang.T("screen_settings")
	}
	return ""
}

// Init tea program
func Run() error {
	program := tea.NewProgram(NewModel())
	_, err := program.Run()
	return err
}
