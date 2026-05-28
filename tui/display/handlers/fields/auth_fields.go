package fields

import "gowordle.com/display/model"

func NextAuthField(f model.AuthField) model.AuthField {
    if f == model.AuthFieldLogin {
        return model.AuthFieldPassword
    }
    return model.AuthFieldLogin
}

func AppendAuthValue(m model.State, value string) model.State {
    if value == "" {
        return m
    }
    switch m.Auth.Field {
    case model.AuthFieldPassword:
        m.Auth.Password += value
    default:
        m.Auth.Login += value
    }
    return m
}

func DeleteAuthValue(m model.State) model.State {
    switch m.Auth.Field {
    case model.AuthFieldPassword:
        if len(m.Auth.Password) > 0 {
            m.Auth.Password = m.Auth.Password[:len(m.Auth.Password)-1]
        }
    default:
        if len(m.Auth.Login) > 0 {
            m.Auth.Login = m.Auth.Login[:len(m.Auth.Login)-1]
        }
    }
    return m
}
