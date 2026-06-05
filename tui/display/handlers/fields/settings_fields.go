package fields

import "gowordle.com/display/model"

var themeOptions = []string{"sombre", "clair"}
var langOptions = []string{"fr", "en"}
var displayOptions = []string{"compact", "normal"}

func NextSettingsField(f model.SettingsField) model.SettingsField {
	return (f + 1) % model.SettingsFieldCount
}

func PrevSettingsField(f model.SettingsField) model.SettingsField {
	if f == 0 {
		return model.SettingsFieldCount - 1
	}
	return f - 1
}

func CycleSettingsValue(s model.Settings) model.Settings {
	switch s.Field {
	case model.SettingsFieldTheme:
		s.Theme = cycleOption(s.Theme, themeOptions)
	case model.SettingsFieldLanguage:
		s.Language = cycleOption(s.Language, langOptions)
	case model.SettingsFieldDisplayMode:
		s.DisplayMode = cycleOption(s.DisplayMode, displayOptions)
	}
	return s
}

func cycleOption(current string, options []string) string {
	for i, v := range options {
		if v == current {
			return options[(i+1)%len(options)]
		}
	}
	return options[0]
}
