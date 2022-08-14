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

	s.router.Get("/rect.png", handlerStatic("rect.png"))
	s.router.Get("/dot.png", handlerStatic("dot.png"))
	s.router.Get("/diamond.png", handlerStatic("diamond.png"))
	s.router.Get("/rhombus.png", handlerStatic("rhombus.png"))
}
