package main

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Localizer struct {
	translations map[string]map[string]string
}

func NewLocalizer(dir string) (*Localizer, error) {
	result := make(map[string]map[string]string)

	files, err := filepath.Glob(dir + "/*.toml")
	if err != nil {
		return nil, err
	}

	for _, path := range files {
		var translations map[string]string

		_, err := toml.DecodeFile(path, &translations)
		if err != nil {
			return nil, err
		}

		// get filename without extension
		locale := filepath.Base(path)
		locale = locale[:len(locale)-len(filepath.Ext(locale))]
		result[locale] = translations
	}

	return &Localizer{translations: result}, nil
}

func (l *Localizer) GetTranslationsForLang(lang string) map[string]string {
	if translations, ok := l.translations[lang]; ok {
		return translations
	}
	return nil
}

func (l *Localizer) HasLanguage(lang string) bool {
	_, ok := l.translations[lang]
	return ok
}

func (l *Localizer) I18n(lang, key string) string {
	if translations, ok := l.translations[lang]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	return key
}
