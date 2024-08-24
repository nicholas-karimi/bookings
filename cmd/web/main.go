package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/nicholas-karimi/bookings/pkg/config"
	"github.com/nicholas-karimi/bookings/pkg/handlers"
	"github.com/nicholas-karimi/bookings/pkg/render"
)

const portNumber = "8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	addr := flag.String("addr", ":"+portNumber, "Serving Http connection")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// togle truw in prod
	app.Inproduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteDefaultMode
	session.Cookie.Secure = app.Inproduction

	app.Session = session

	template_cache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}

	app.TemplateCache = template_cache
	// app.UseCache = false

	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:     *addr,
		Handler:  routes(&app),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
