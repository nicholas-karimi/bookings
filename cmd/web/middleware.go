package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get url user visiting
		// path, _ := r.URL()
		app.InfoLog.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)

	})
}

// csrf

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
