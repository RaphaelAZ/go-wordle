package display

import (
	"strings"

	"charm.land/lipgloss/v2"
	"gowordle.com/display/handlers"
	"gowordle.com/display/model"
)

func HomeScreen() string {
	theme := DefaultTheme()
	body := []string{
		theme.Section.Render("[1] Accueil"),
		theme.Section.Render("[2] Connexion"),
		theme.Section.Render("[3] Paramètres"),
	}

	return RenderApp(
		"Go Wordle",
		"Interface terminal du projet",
		body,
		theme.Accent.Render("Navigation au clavier à venir"),
	)
}

func AuthScreen(m model.State) string {
	theme := DefaultTheme()
	body := []string{
		authFieldLine(theme, "Login", m.Auth.Login, m.Editing && m.Auth.Field == model.AuthFieldLogin, false),
		authFieldLine(theme, "Mot de passe", m.Auth.Password, m.Editing && m.Auth.Field == model.AuthFieldPassword, true),
	}

	footer := "Appuyez sur Entrée pour éditer"
	if m.Editing {
		footer = "Entrée: champ suivant, Échap: quitter l'édition"
	}

	return RenderApp(
		"Authentification",
		"Connexion au compte utilisateur",
		body,
		theme.Muted.Render(footer),
	)
}

func SettingsScreen() string {
	theme := DefaultTheme()
	body := []string{
		theme.Text.Render("Thème          : clair / sombre"),
		theme.Text.Render("Langue         : fr"),
		theme.Text.Render("Mode d'affichage: terminal compact"),
	}

	return RenderApp(
		"Paramètres",
		"Configuration visuelle de l'application",
		body,
		theme.Muted.Render("Prévu pour être branché sur un stockage local JSON"),
	)
}

func GameScreen(m model.State) string {
	if m.Game.WordToGuess == "" {
		m.Game.WordToGuess = "BRUME"
	}

	width := m.Width
	if width <= 0 {
		width = 120
	}
	height := m.Height
	if height <= 0 {
		height = 30
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

	footer := "Tapez une lettre, Entrée pour valider, Échap pour le menu"
	if m.Game.Status == model.GameWon {
		footer = "Bravo ! Vous avez trouvé le mot !"
	} else if m.Game.Status == model.GameLost {
		footer = "Perdu ! Le mot était : " + strings.ToUpper(m.Game.WordToGuess)
	}

	grid := lipgloss.NewStyle().Padding(1, 2).Render(strings.Join(gridLines, "\n"))
	keyboard := lipgloss.NewStyle().Padding(1, 2).Render(renderKeyboard(m, theme))
	middle := lipgloss.Place(width, height-4,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, grid, "    ", keyboard),
	)

	title := theme.Title.Width(width).Render("Wordle")
	footerLine := theme.Muted.Width(width).Align(lipgloss.Center).Render(footer)

	return lipgloss.JoinVertical(lipgloss.Left, title, middle, footerLine)
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
		SettingsScreen(),
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
