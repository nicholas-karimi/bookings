package render

import (
	"net/http"
	"testing"

	"github.com/nicholas-karimi/bookings/internals/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// test data
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value not found in session")
	}

	if result == nil {
		t.Error("AddDefaultData returned a nil value")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplates(&ww,"home.page.tmpl", r, &models.TemplateData{})

	if err != nil {
		t.Error("error writing template to browser")
	}
	
	err = RenderTemplates(&ww,  "none-existent.page.tmpl", r, &models.TemplateData{})

	if err != nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession()(*http.Request, error){
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}


func TestNewTemplates(t *testing.T){
	NewTemplates(app)
}


func TestCreateTemplateCache(t *testing.T){
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}