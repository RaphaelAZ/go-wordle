package model

type AuthField int

type Auth struct {
	Field    AuthField
	Login    string
	Password string
	Error    string
	Loading  bool
}

// Auth fields
const (
	AuthFieldLogin AuthField = iota
	AuthFieldPassword
)

type AuthResultMsg struct {
	Token string
	Err   error
}
