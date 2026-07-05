package fields

import "gowordle.com/display/model"

// authFieldIndex returns the position of f within the fields of the given mode.
func authFieldIndex(fs []model.AuthField, f model.AuthField) int {
    for i, cur := range fs {
        if cur == f {
            return i
        }
    }
    return 0
}

func NextAuthField(a model.Auth) model.AuthField {
    fs := a.Fields()
    i := authFieldIndex(fs, a.Field)
    return fs[(i+1)%len(fs)]
}

// IsLastAuthField reports whether the focused field is the last one, i.e. the
// submit field for the current mode.
func IsLastAuthField(a model.Auth) bool {
    fs := a.Fields()
    return authFieldIndex(fs, a.Field) == len(fs)-1
}

// ToggleAuthMode switches between login and registration, resetting the focus
// to the first field and clearing any previous error.
func ToggleAuthMode(a model.Auth) model.Auth {
    if a.Mode == model.AuthModeRegister {
        a.Mode = model.AuthModeLogin
    } else {
        a.Mode = model.AuthModeRegister
    }
    a.Field = a.Fields()[0]
    a.Error = ""
    return a
}

func AppendAuthValue(m model.State, value string) model.State {
    if value == "" {
        return m
    }
    switch m.Auth.Field {
    case model.AuthFieldUsername:
        m.Auth.Username += value
    case model.AuthFieldPassword:
        m.Auth.Password += value
    default:
        m.Auth.Email += value
    }
    return m
}

func DeleteAuthValue(m model.State) model.State {
    switch m.Auth.Field {
    case model.AuthFieldUsername:
        if len(m.Auth.Username) > 0 {
            m.Auth.Username = m.Auth.Username[:len(m.Auth.Username)-1]
        }
    case model.AuthFieldPassword:
        if len(m.Auth.Password) > 0 {
            m.Auth.Password = m.Auth.Password[:len(m.Auth.Password)-1]
        }
    default:
        if len(m.Auth.Email) > 0 {
            m.Auth.Email = m.Auth.Email[:len(m.Auth.Email)-1]
        }
    }
    return m
}
