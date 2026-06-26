package display

import (
	"strings"

	"charm.land/lipgloss/v2"
	"gowordle.com/display/handlers"
	"gowordle.com/display/lang"
	"gowordle.com/display/model"
)

func HomeScreen() string {
	theme := DefaultTheme()
	body := []string{
		theme.Section.Render("[1] " + lang.T("screen_home")),
		theme.Section.Render("[2] " + lang.T("screen_auth")),
		theme.Section.Render("[3] " + lang.T("screen_settings")),
	}

	return RenderApp(
		lang.T("home_title"),
		lang.T("home_subtitle"),
		body,
		"",
	)
}

func AuthScreen(m model.State) string {
	theme := DefaultTheme()
	body := []string{
		authFieldLine(theme, lang.T("auth_field_email"), m.Auth.Login, m.Editing && m.Auth.Field == model.AuthFieldLogin, false),
		authFieldLine(theme, lang.T("auth_field_password"), m.Auth.Password, m.Editing && m.Auth.Field == model.AuthFieldPassword, true),
	}

	if m.Auth.Error != "" {
		body = append(body, theme.Error.Render(lang.T("auth_error", map[string]any{"Error": m.Auth.Error})))
	}

	footer := lang.T("auth_footer_idle")
	switch {
	case m.Auth.Loading:
		footer = lang.T("auth_footer_loading")
	case m.Editing:
		footer = lang.T("auth_footer_editing")
	}

	return RenderApp(
		lang.T("auth_title"),
		lang.T("auth_subtitle"),
		body,
		theme.Muted.Render(footer),
	)
}

func SettingsScreen(m model.State) string {
	theme := DefaultTheme()

	type field struct {
		label    string
		value    string
		id       model.SettingsField
		isAction bool
	}
	fieldDefs := []field{
		{label: lang.T("settings_field_theme"), value: lang.T("settings_theme_" + m.Settings.Theme), id: model.SettingsFieldTheme},
		{label: lang.T("settings_field_language"), value: m.Settings.Language, id: model.SettingsFieldLanguage},
		{label: lang.T("settings_field_display"), value: m.Settings.DisplayMode, id: model.SettingsFieldDisplayMode},
	}
	if m.Connected {
		fieldDefs = append(fieldDefs, field{label: lang.T("settings_field_logout"), id: model.SettingsFieldLogout, isAction: true})
	}

	var body []string
	for _, f := range fieldDefs {
		active := m.Editing && m.Settings.Field == f.id
		style := theme.Text
		prefix := "  "
		if active {
			style = theme.Accent
			prefix = "> "
		}
		var line string
		if f.isAction {
			hint := ""
			if active {
				hint = "  ↵"
			}
			line = style.Render(prefix + f.label + hint)
		} else {
			hint := ""
			if active {
				hint = "  ← →"
			}
			line = style.Render(prefix + f.label + " : " + f.value + hint)
		}
		body = append(body, line)
	}

	footer := lang.T("settings_footer_idle")
	if m.Editing {
		footer = lang.T("settings_footer_editing")
	}

	return RenderApp(
		lang.T("settings_title"),
		lang.T("settings_subtitle"),
		body,
		theme.Muted.Render(footer),
	)
}

func GameLoadingScreen(m model.State) string {
	theme := DefaultTheme()
	msg := theme.Muted.Render(lang.T("game_loading"))
	if m.Settings.DisplayMode == "compact" {
		return msg
	}
	width := m.Width
	if width <= 0 {
		width = 120
	}
	height := m.Height
	if height <= 0 {
		height = 30
	}
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, msg)
}

func GameScreen(m model.State) string {
	if m.Game.WordLoading {
		return GameLoadingScreen(m)
	}

	if m.Game.WordError != "" {
		theme := DefaultTheme()
		msg := theme.Error.Render(lang.T("game_word_error", map[string]any{"Error": m.Game.WordError}))
		if m.Settings.DisplayMode == "compact" {
			return RenderApp(lang.T("game_title"), "", []string{msg}, "")
		}
		width := m.Width
		if width <= 0 {
			width = 120
		}
		height := m.Height
		if height <= 0 {
			height = 30
		}
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, msg)
	}

	theme := DefaultTheme()

	var gridLines []string
	for i := 0; i < 6; i++ {
		if i == len(m.Game.TriedWords) {
			padded := handlers.AddPaddingsToCurrentGuessedWord(m.Game.CurrentGuess)
			gridLines = append(gridLines, theme.Accent.Render(handlers.ToWordleDisplayString(padded)))
			continue
		}
		if i > len(m.Game.TriedWords) {
			gridLines = append(gridLines, theme.Muted.Render("_ _ _ _ _"))
			continue
		}
		var letters []string
		for _, ls := range m.Game.TriedWords[i] {
			letters = append(letters, handlers.RenderWordleColoredLetter(ls, theme.Letters.Correct, theme.Letters.Misplaced, theme.Letters.Incorrect))
		}
		gridLines = append(gridLines, strings.Join(letters, " "))
	}

	footer := lang.T("game_footer_playing")
	if m.Game.Status == model.GameWon {
		footer = lang.T("game_footer_won")
	} else if m.Game.Status == model.GameLost {
		footer = lang.T("game_footer_lost", map[string]any{"Word": strings.ToUpper(m.Game.WordToGuess)})
	}

	if m.Settings.DisplayMode == "compact" {
		return compactGameScreen(m, theme, gridLines, footer)
	}

	width := m.Width
	if width <= 0 {
		width = 120
	}
	height := m.Height
	if height <= 0 {
		height = 30
	}

	grid := lipgloss.NewStyle().Padding(1, 2).Render(strings.Join(gridLines, "\n"))
	keyboard := lipgloss.NewStyle().Padding(1, 2).Render(renderKeyboard(m, theme))
	middle := lipgloss.Place(width, height-4,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, grid, "    ", keyboard),
	)

	title := theme.Title.Width(width).Render(lang.T("game_title"))
	footerLine := theme.Muted.Width(width).Align(lipgloss.Center).Render(footer)

	lines := []string{title, middle, footerLine}
	if m.Game.Status != model.GamePlaying {
		if m.Game.SaveLoading {
			lines = append(lines, theme.Muted.Width(width).Align(lipgloss.Center).Render(lang.T("game_saving")))
		} else if m.Game.SaveError != "" {
			lines = append(lines, theme.Error.Width(width).Align(lipgloss.Center).Render(lang.T("game_save_error", map[string]any{"Error": m.Game.SaveError})))
		}
		lines = append(lines, theme.Button.Width(width).Align(lipgloss.Center).Render(lang.T("game_restart_button")))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func compactGameScreen(m model.State, theme Theme, gridLines []string, footer string) string {
	grid := lipgloss.NewStyle().Padding(0, 2).Render(strings.Join(gridLines, "\n"))
	kb := lipgloss.NewStyle().Padding(0, 1).Render(renderKeyboard(m, theme))
	body := []string{
		lipgloss.JoinHorizontal(lipgloss.Center, grid, "  ", kb),
	}
	if m.Game.Status != model.GamePlaying {
		if m.Game.SaveLoading {
			body = append(body, theme.Muted.Render(lang.T("game_saving")))
		} else if m.Game.SaveError != "" {
			body = append(body, theme.Error.Render(lang.T("game_save_error", map[string]any{"Error": m.Game.SaveError})))
		}
		body = append(body, theme.Button.Render(lang.T("game_restart_button")))
	}
	return RenderApp(lang.T("game_title"), "", body, theme.Muted.Render(footer))
}

func renderKeyboard(m model.State, theme Theme) string {
	rows := [][]string{
		{"A", "Z", "E", "R", "T", "Y", "U", "I", "O", "P"},
		{"Q", "S", "D", "F", "G", "H", "J", "K", "L", "M"},
		{"W", "X", "C", "V", "B", "N"},
	}

	var renderedRows []string
	for _, row := range rows {
		var keys []string
		for _, key := range row {
			letter := key[0]
			if ls, ok := m.Game.UsedLetters[letter]; ok {
				keys = append(keys, handlers.RenderWordleColoredLetter(ls, theme.Letters.Correct, theme.Letters.Misplaced, theme.Letters.Incorrect))
			} else {
				keys = append(keys, theme.Text.Render(key))
			}
		}
		renderedRows = append(renderedRows, strings.Join(keys, " "))
	}
	return strings.Join(renderedRows, "\n")
}

func Dashboard() string {
	theme := DefaultTheme()
	sections := []string{
		theme.Button.Render("CUI CUI JE SUIS UN POULET"),
		"",
		HomeScreen(),
		"",
		AuthScreen(model.State{}),
		"",
		GameScreen(model.State{}),
		"",
		SettingsScreen(model.State{}),
	}

	return strings.Join(sections, "\n\n")
}

func authFieldLine(theme Theme, label, value string, active, masked bool) string {
	displayValue := value
	if masked {
		displayValue = strings.Repeat("*", len(value))
	}
	if active {
		displayValue += "_"
	}
	if displayValue == "" {
		displayValue = "_"
	}

	style := theme.Text
	if active {
		style = theme.Accent
	}

	return style.Render(label + " : " + displayValue)
}
