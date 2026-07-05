package model

type AuthField int

type AuthMode int

type Auth struct {
	Mode     AuthMode
	Field    AuthField
	Username string
	Email    string
	Password string
	Error    string
	Loading  bool
}

// Auth modes
const (
	AuthModeLogin AuthMode = iota
	AuthModeRegister
)

// Auth fields
const (
	AuthFieldUsername AuthField = iota
	AuthFieldEmail
	AuthFieldPassword
)

// AuthFields returns the ordered fields shown for the current mode.
// Login only needs email and password; registration also asks for a username.
func (a Auth) Fields() []AuthField {
	if a.Mode == AuthModeRegister {
		return []AuthField{AuthFieldUsername, AuthFieldEmail, AuthFieldPassword}
	}
	return []AuthField{AuthFieldEmail, AuthFieldPassword}
}

type AuthResultMsg struct {
	Token string
	Err   error
}

type WordResultMsg struct {
	WordID int
	Word   string
	Err    error
}

type GameSaveResultMsg struct {
	Err error
}
