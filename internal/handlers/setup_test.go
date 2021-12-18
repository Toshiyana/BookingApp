package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Toshiyana/BookingApp/internal/config"
	"github.com/Toshiyana/BookingApp/internal/models"
	"github.com/Toshiyana/BookingApp/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func getRoutes() http.Handler {
	//----------------------------------------------------------------
	// run() in main.go
	//----------------------------------------------------------------
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// use session in handlers
	session = scs.New()
	session.Lifetime = 24 * time.Hour // last for 24 hours
	// set cookie parameters because all session use cookies in one form.
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	// In develop mode, Usecache sets false because of reloding templates. <- check templates changed.
	// In release mode, Usecashe sets true because of not reloding templates.
	app.UseCache = true // In test, not to call CreateTemplateCache() and use pathToTemplates in render.go, set true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)
	//----------------------------------------------------------------
	// routes() in routes.go
	//----------------------------------------------------------------
	// Use chi library for router
	mux := chi.NewRouter()

	// Middleware allows you to process a request as it comes into your Web application
	// and perform some action on it.
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf) // turn on a middleware
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(Repo.Home))
	mux.Get("/about", http.HandlerFunc(Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(Repo.Majors))

	mux.Get("/make-reservation", http.HandlerFunc(Repo.Reservation))
	mux.Post("/make-reservation", http.HandlerFunc(Repo.PostReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(Repo.ReservationSummary))

	mux.Get("/search-availability", http.HandlerFunc(Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(Repo.AvailabilityJSON))

	mux.Get("/contact", http.HandlerFunc(Repo.Contact))

	// create a file server, a place to get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Set base cookie because it uses cookies to make sure that the token it generates is available on a per page basis.
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",              // apply to the entire site for cookie secure
		Secure:   app.InProduction, // In production, change true
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a template
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	// get the template cache from the app config

	// create the template cache only one time

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		fmt.Println("Page is currently", page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
