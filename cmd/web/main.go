package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Toshiyana/BookingApp/driver"
	"github.com/Toshiyana/BookingApp/internal/config"
	"github.com/Toshiyana/BookingApp/internal/handlers"
	"github.com/Toshiyana/BookingApp/internal/helpers"
	"github.com/Toshiyana/BookingApp/internal/models"
	"github.com/Toshiyana/BookingApp/internal/render"
	"github.com/joho/godotenv"

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

	// fmt.Printf("Starting application on port %s\n", portNumber)
	fmt.Printf("Server starting: http://localhost%s\n", portNumber)

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
	gob.Register(map[string]int{})

	// Env Parameters (You can use .env or Read flags)
	//----------------------------------------------------------------
	// .env and gotodoenv library
	//----------------------------------------------------------------
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("cannot read .env file")
		os.Exit(1)
	}

	inProduction, _ := strconv.ParseBool(os.Getenv("PRODUCTION"))
	useCache, _ := strconv.ParseBool(os.Getenv("CACHE"))
	dbHost := os.Getenv("POSTGRESQL_HOST")
	dbName := os.Getenv("POSTGRESQL_DBNAME")
	dbUser := os.Getenv("POSTGRESQL_USERNAME")
	dbPass := os.Getenv("POSTGRESQL_PASSWORD")
	dbPort := os.Getenv("POSTGRESQL_PORT")
	dbSSL := os.Getenv("POSTGRESQL_SSL")

	//----------------------------------------------------------------
	// Read flags
	//----------------------------------------------------------------
	// Second parameter is default value.
	// inProduction := flag.Bool("production", true, "Application is in production")// default is "true" because of safety
	// useCache := flag.Bool("cache", true, "Use template cache")
	// dbHost := flag.String("dbhost", "localhost", "Database host")// default is "localhost" because you rarely have the database on the same machine as your application server.
	// dbName := flag.String("dbname", "", "Database name")
	// dbUser := flag.String("dbuser", "", "Database user")
	// dbPass := flag.String("dbpass", "", "Database password")
	// dbPort := flag.String("dbport", "5432", "Database port")
	// dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	// flag.Parse()

	// if *dbName == "" || *dbUser == "" {
	// 	fmt.Println("Missing required flags")
	// 	os.Exit(1)
	// }

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = inProduction
	// In develop mode, Usecache sets false because of reloding templates. <- check templates changed.
	// In release mode, Usecashe sets true because of not reloding templates.
	app.UseCache = useCache

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
	// connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)
	db, err := driver.ConnectSQL(connectionString)
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

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}
