package model

// Navigation state
type State struct {
	Selected Screen
	Width    int
	Height   int
	Editing  bool
	Focus    bool
	Auth     Auth
}
