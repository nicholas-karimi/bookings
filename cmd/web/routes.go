package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nicholas-karimi/bookings/pkg/config"
	"github.com/nicholas-karimi/bookings/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := http.NewServeMux()

	// mux.HandleFunc("/", handlers.Repo.Home)
	// mux.HandleFunc("/about", handlers.Repo.About)
	mux := chi.NewRouter()
	// mux.Use(NoSurf)
	// mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	// return NoSurf(WriteToConsole(mux))
	return mux
}
