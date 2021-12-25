package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Toshiyana/BookingApp/driver"
	"github.com/Toshiyana/BookingApp/internal/config"
	"github.com/Toshiyana/BookingApp/internal/handlers"
	"github.com/Toshiyana/BookingApp/internal/helpers"
	"github.com/Toshiyana/BookingApp/internal/models"
	"github.com/Toshiyana/BookingApp/internal/render"

	"github.com/alexedwards/scs/v2"
)

// declare portNumber out of main() because I never want the portNumber to be changed by another part of the application.
const portNumber = ":8080"

var app config.AppConfig // also use in Nosurf() of middleware.go
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget close
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	// Use routes
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.Reservation{})
	gob.Register(models.RoomRestriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// use session in handlers
	session = scs.New()
	session.Lifetime = 24 * time.Hour // last for 24 hours
	// set cookie parameters because all session use cookies in one form.
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go-bookingapp user=toshi password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	// In develop mode, Usecache sets false because of reloding templates. <- check templates changed.
	// In release mode, Usecashe sets true because of not reloding templates.
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}
