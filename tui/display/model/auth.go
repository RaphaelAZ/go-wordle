package model

type AuthField int

type Auth struct {
	Field    AuthField
	Login    string
	Password string
}

// Auth fields
const (
	AuthFieldLogin AuthField = iota // Increment values for each screen
	AuthFieldPassword
)
