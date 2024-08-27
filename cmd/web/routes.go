package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nicholas-karimi/bookings/internals/config"
	"github.com/nicholas-karimi/bookings/internals/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	mux.Use(NoSurf)
	// mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/index", handlers.Repo.IndexPage)
	mux.Get("/about", handlers.Repo.About)

	mux.Get("/generals-quarters", handlers.Repo.GeneralsPage)
	mux.Get("/majors-suite", handlers.Repo.MajorsPage)

	mux.Get("/search-availability", handlers.Repo.AvailabilityPage)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-avalibility-json", handlers.Repo.AvailabilityJsonData)

	mux.Get("/make-reservations", handlers.Repo.MakeReservationPage)
	mux.Post("/make-reservations", handlers.Repo.PostReservationPage)

	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)


	mux.Get("/contact", handlers.Repo.ContactPage)

	// serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// return NoSurf(WriteToConsole(mux))
	return mux
}
