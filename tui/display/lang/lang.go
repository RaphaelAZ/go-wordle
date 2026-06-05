package lang

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales
var localesFS embed.FS

var loc *i18n.Localizer

func Init(langCode string) {
	tag := language.French
	if langCode == "en" {
		tag = language.English
	}

	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, _ = bundle.LoadMessageFileFS(localesFS, "locales/active.fr.toml")
	_, _ = bundle.LoadMessageFileFS(localesFS, "locales/active.en.toml")

	loc = i18n.NewLocalizer(bundle, tag.String())
}

func T(id string, data ...map[string]any) string {
	if loc == nil {
		Init("fr")
	}
	cfg := &i18n.LocalizeConfig{MessageID: id}
	if len(data) > 0 {
		cfg.TemplateData = data[0]
	}
	str, err := loc.Localize(cfg)
	if err != nil {
		return id
	}
	return str
}
