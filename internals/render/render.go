package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/nicholas-karimi/bookings/internals/config"
	"github.com/nicholas-karimi/bookings/internals/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(template_data *models.TemplateData, r *http.Request) *models.TemplateData {
	template_data.CSRFToken = nosurf.Token(r)
	return template_data
}

func RenderTemplates(w http.ResponseWriter, tmpl string, r *http.Request, template_data *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		// get template from cache app cinfig
		tc = app.TemplateCache

	} else {
		// read from disk
		tc, _ = CreateTemplateCache()
	}

	tmplate, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	template_data = AddDefaultData(template_data, r)
	_ = tmplate.Execute(buf, template_data)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
	}
}

// create template cache as map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
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

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}

		}
		myCache[name] = ts
	}
	return myCache, nil
}
