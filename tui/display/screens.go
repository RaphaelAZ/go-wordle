package display

import (
	"strings"

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
	theme := DefaultTheme()
	wordToGuess := m.Game.WordToGuess
	if wordToGuess == "" {
		wordToGuess = "_ _ _ _ _"
	}

	triedWords := "Aucun mot essayé pour le moment"
	if len(m.Game.TriedWords) > 0 {
		triedWords = strings.Join(m.Game.TriedWords, "  |  ")
	}

	body := []string{
		theme.Text.Render("Mot à deviner    : " + wordToGuess),
		theme.Text.Render("Mots essayés     : " + triedWords),
		theme.Text.Render("Essais restants  : 6"),
	}

	return RenderApp(
		"Jeu",
		"Wordle en cours",
		body,
		theme.Muted.Render("Prochaine étape: brancher la logique de proposition et de validation"),
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
