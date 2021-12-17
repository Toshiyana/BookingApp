package main

import (
	"net/http"

	"github.com/Toshiyana/BookingApp/pkg/config"
	"github.com/Toshiyana/BookingApp/pkg/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// Use chi library for router
	mux := chi.NewRouter()

	// Middleware allows you to process a request as it comes into your Web application
	// and perform some action on it.
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf) // turn on a middleware
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.Majors))
	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.Reservation))

	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostAvailability))

	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))

	// create a file server, a place to get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
