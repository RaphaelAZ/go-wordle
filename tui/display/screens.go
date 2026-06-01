package display

import (
	"strings"

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

	theme := DefaultTheme()
	body := []string{}

	for i := 0; i < 6; i++ {
		// if current word
		if i == len(m.Game.TriedWords) {
			// No styling at all
			padded := handlers.AddPaddingsToCurrentGuessedWord(m.Game.CurrentGuess)
			body = append(body, theme.Accent.Render(handlers.ToWordleDisplayString(padded)))
			continue
		}

		if i > len(m.Game.TriedWords) {
			body = append(body, theme.Muted.Render("_ _ _ _ _"))
			continue
		}

		var toAppend []string

		// render the letters of the game
		for j := 0; j < len(m.Game.TriedWords[i]); j++ {
			toAppend = append(toAppend, handlers.RenderWordleColoredLetter(m.Game.TriedWords[i][j], theme.Letters.Correct, theme.Letters.Misplaced, theme.Letters.Incorrect))
		}

		body = append(body, strings.Join(toAppend, " "))
	}

	footer := "Tapez une lettre, Entrée pour valider, Échap pour le menu"
	if m.Game.Status == model.GameWon {
		footer = "Bravo ! Vous avez trouvé le mot !"
	} else if m.Game.Status == model.GameLost {
		footer = "Perdu ! Le mot était : " + strings.ToUpper(m.Game.WordToGuess)
	}

	return RenderApp(
		"Jeu",
		"Wordle en cours",
		body,
		theme.Muted.Render(footer),
	)
}

func Dashboard() string {
	theme := DefaultTheme()
	sections := []string{
		theme.Button.Render("CLI TUI"),
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
