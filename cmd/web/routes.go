package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nicholas-karimi/bookings/pkg/config"
	"github.com/nicholas-karimi/bookings/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	// mux.Use(NoSurf)
	// mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	// serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// return NoSurf(WriteToConsole(mux))
	return mux
}
