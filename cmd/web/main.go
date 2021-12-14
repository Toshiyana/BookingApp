package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Toshiyana/BookingApp/pkg/config"
	"github.com/Toshiyana/BookingApp/pkg/handlers"
	"github.com/Toshiyana/BookingApp/pkg/render"

	"github.com/alexedwards/scs/v2"
)

// declare portNumber out of main() because I never want the portNumber to be changed by another part of the application.
const portNumber = ":8080"

var app config.AppConfig // also use in Nosurf() of middleware.go
var session *scs.SessionManager

// main is the main application function
func main() {
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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	// In develop mode, Usecache sets false because of reloding templates. <- check templates changed.
	// In release mode, Usecashe sets true because of not reloding templates.
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	// Use routes
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	// http.HandleFunc("/divide", Divide)

	// start web server that listens for requests
	// listen to particular port of my computer (In this case, listen on 8080)
	// _ = http.ListenAndServe(portNumber, nil)
	// access localhost:8080 in browser
}

// command in terminal: go run cmd/web/*.go
