package state

// Screens
type Screen int

const (
	ScreenHome Screen = iota
	ScreenAuth
	ScreenSettings
)

type AuthField int

const (
	AuthFieldLogin AuthField = iota
	AuthFieldPassword
)

var ScreenLabels = []string{
	"Accueil",
	"Connexion",
	"Paramètres",
}

// Navigation state
type Model struct {
	Selected  Screen
	Width     int
	Height    int
	Editing   bool
	Focus     bool
	AuthField AuthField
	Login     string
	Password  string
}
