package main

import (
	"encoding/gob"
	"flag"
	"github.com/nicholas-karimi/bookings/internals/driver"
	"github.com/nicholas-karimi/bookings/internals/helpers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/nicholas-karimi/bookings/internals/config"
	"github.com/nicholas-karimi/bookings/internals/handlers"
	"github.com/nicholas-karimi/bookings/internals/models"
	"github.com/nicholas-karimi/bookings/internals/render"
)

const portNumber = "8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	addr := flag.String("addr", ":"+portNumber, "Serving Http connection")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	srv := &http.Server{
		Addr:     *addr,
		Handler:  routes(&app),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func run() (*driver.DB, error) {
	// store session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	// togle truw in prod
	app.Inproduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteDefaultMode
	session.Cookie.Secure = app.Inproduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectionSQL("host=localhost port=5432 dbname=bookings user=admin password=Incorrect")

	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	log.Println("Connected to database.")

	template_cache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = template_cache
	// app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
