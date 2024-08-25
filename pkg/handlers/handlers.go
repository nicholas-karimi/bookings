package handlers

import (
	"net/http"

	"github.com/nicholas-karimi/bookings/pkg/config"
	"github.com/nicholas-karimi/bookings/pkg/models"
	"github.com/nicholas-karimi/bookings/pkg/render"
)

// setups the repository for the handlers
var Repo *Repository

// repository type
type Repository struct {
	App *config.AppConfig
}

// cretaes a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// sets repo forthe handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplates(w, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := repo.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplates(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) IndexPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "index.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) GeneralsPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "generals.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) MajorsPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "majors.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) ReservationPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "reservation.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) MakeReservationPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "make-reservations.page.tmpl", &models.TemplateData{})
}
func (repo *Repository) ContactPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "contact.page.tmpl", &models.TemplateData{})
}
