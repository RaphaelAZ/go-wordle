package model

type Screen int

// Available screens
const (
	ScreenHome Screen = iota // Increment values for each screen
	ScreenAuth
	ScreenSettings
)

// Screens labels
var ScreenLabels = [3]string{
	"Accueil",
	"Connexion",
	"Paramètres",
}
