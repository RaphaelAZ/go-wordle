package model

type SettingsField int

const (
	SettingsFieldTheme SettingsField = iota
	SettingsFieldLanguage
	SettingsFieldDisplayMode
	SettingsFieldCount
)

type Settings struct {
	Field       SettingsField
	Theme       string // "sombre" | "clair"
	Language    string // "fr" | "en"
	DisplayMode string // "compact" | "normal"
}

type SettingsChangedMsg struct{}

func DefaultSettings() Settings {
	return Settings{
		Theme:       "sombre",
		Language:    "fr",
		DisplayMode: "compact",
	}
}
