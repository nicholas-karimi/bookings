package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
	"github.com/nicholas-karimi/bookings/internals/config"
	"github.com/nicholas-karimi/bookings/internals/models"
	"github.com/nicholas-karimi/bookings/internals/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// store sessiom
	gob.Register(models.Reservation{})

	// togle truw in prod
	app.Inproduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteDefaultMode
	session.Cookie.Secure = app.Inproduction

	app.Session = session

	template_cache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}

	app.TemplateCache = template_cache
	app.UseCache = true

	render.NewRenderer(&app)

	repo := NewRepo(&app)
	NewHandlers(repo)

	mux := chi.NewRouter()
	// mux.Use(NoSurf)
	// mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/index", Repo.IndexPage)
	mux.Get("/about", Repo.About)

	mux.Get("/generals-quarters", Repo.GeneralsPage)
	mux.Get("/majors-suite", Repo.MajorsPage)

	mux.Get("/search-availability", Repo.AvailabilityPage)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-avalibility-json", Repo.AvailabilityJsonData)

	mux.Get("/make-reservations", Repo.MakeReservationPage)
	mux.Post("/make-reservations", Repo.PostReservationPage)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.ContactPage)

	// serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// return NoSurf(WriteToConsole(mux))
	return mux

}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.Inproduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)

}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// log.Println("page is currently being parsed: ", page)
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
