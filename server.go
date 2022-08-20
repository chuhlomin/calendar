package main

import (
	"io/ioutil"
	"log"
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

	files, err := ioutil.ReadDir("./static")
	if err != nil {
		log.Fatalf("could not read static files: %v", err)
	}
	for _, file := range files {
		s.router.Get("/"+file.Name(), handlerStatic(file.Name()))
	}

	s.router.Post("/pdf", handlerPDF)
}
