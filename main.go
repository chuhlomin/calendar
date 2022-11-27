package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var localizer *Localizer

func main() {
	log.Println("Starting...")

	if err := run(); err != nil {
		log.Printf("ERROR %v", err)
	}

	log.Println("Stopped")
}

func run() error {
	bind := flag.String("bind", ":8080", "Bind address")
	flag.Parse()

	var err error
	localizer, err = NewLocalizer("langs")
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5, "application/json"))

	srv := server{router: r}
	srv.routes()

	log.Printf("Starting server on %s", *bind)
	if err := http.ListenAndServe(*bind, srv.router); err != nil {
		return err
	}

	return nil
}
