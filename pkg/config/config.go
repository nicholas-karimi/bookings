package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// Holds appwide config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger

	Inproduction bool
	Session      *scs.SessionManager
}
