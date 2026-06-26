package model

type SettingsField int

const (
	SettingsFieldTheme SettingsField = iota
	SettingsFieldLanguage
	SettingsFieldDisplayMode
	SettingsFieldLogout
	SettingsFieldCount
)

type Settings struct {
	Field       SettingsField
	Theme       string // "sombre" | "clair"
	Language    string // "fr" | "en"
	DisplayMode string // "compact" | "normal"
	BackendUrl  string // "https://gowordle.alwaysdata.net"
}

type SettingsChangedMsg struct{}
type LogoutMsg struct{}

func DefaultSettings() Settings {
	return Settings{
		Theme:       "sombre",
		Language:    "fr",
		DisplayMode: "compact",
		BackendUrl:  "https://gowordle.alwaysdata.net",
	}
}
