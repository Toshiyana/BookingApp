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
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.Majors))

	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.Reservation))
	mux.Post("/make-reservation", http.HandlerFunc(handlers.Repo.PostReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(handlers.Repo.ReservationSummary))

	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(handlers.Repo.AvailabilityJSON))

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// create a file server, a place to get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth) // use Auth in middleware. In Product mode, uncomment this for security
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)

		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)

		mux.Get("/process-reservation/{src}/{id}", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}", handlers.Repo.AdminDeleteReservation)
	})

	return mux
}
