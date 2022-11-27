package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerI18n(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}

	translations := localizer.GetTranslationsForLang(lang)
	if translations == nil {
		log.Printf("could not find translations for lang %q", lang)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(translations)
	if err != nil {
		log.Printf("error encoding translations: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
