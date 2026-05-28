package model

type Screen int

// Available screens
const (
	ScreenHome Screen = iota // Increment values for each screen
	ScreenAuth
	ScreenGame
	ScreenSettings
)

// Screens labels
var ScreenLabels = [4]string{
	"Accueil",
	"Connexion",
	"Jeu",
	"Paramètres",
}
