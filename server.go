package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	router chi.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.Use(middleware.StripSlashes)

	s.router.Get("/", handlerStatic("static", "index.html"))

	staticFilesDirs := []string{"static", "fonts"}

	for _, dir := range staticFilesDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("could not read static files: %v", err)
		}
		for _, file := range files {
			s.router.Get("/"+file.Name(), handlerStatic(dir, file.Name()))
		}
	}

	s.router.Post("/pdf", handlerPDF)
	s.router.Post("/encode", handlerEncode)
	s.router.Get("/i18n", handlerI18n)
}
