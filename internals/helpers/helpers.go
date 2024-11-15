package helpers

import (
	"fmt"
	"github.com/nicholas-karimi/bookings/internals/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers sets up config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	//write to infolog
	app.InfoLog.Println("Client Error:", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
