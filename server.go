package main

import (
	"net/http"

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

	s.router.Get("/", handlerStatic("index.html"))
	s.router.Get("/styles.css", handlerStatic("styles.css"))
	s.router.Get("/script.js", handlerStatic("script.js"))
	s.router.Get("/patterns.png", handlerStatic("patterns.png"))
	s.router.Get("/GitHub-Mark-32px.png", handlerStatic("GitHub-Mark-32px.png"))
	s.router.Get("/GitHub-Mark-Light-32px.png", handlerStatic("GitHub-Mark-Light-32px.png"))

	s.router.Post("/pdf", handlerPDF)
}
