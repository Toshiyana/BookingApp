package main

import (
	"net/http"

	"github.com/Toshiyana/BookingApp/internal/config"
	"github.com/Toshiyana/BookingApp/internal/handlers"

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

	// You can remove "http.HandlerFunc()"
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// create a file server, a place to get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth) // use Auth in middleware. In Product mode, uncomment this for security
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)

		mux.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)

		mux.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)
	})

	return mux
}
