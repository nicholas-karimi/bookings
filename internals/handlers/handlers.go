package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicholas-karimi/bookings/internals/config"
	"github.com/nicholas-karimi/bookings/internals/models"
	"github.com/nicholas-karimi/bookings/internals/render"
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

	render.RenderTemplates(w, "home.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := repo.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplates(w, "about.page.tmpl", r, &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) IndexPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "index.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) GeneralsPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "generals.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) MajorsPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "majors.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) AvailabilityPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "search-availability.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.PostForm.Get("start")
	end := r.PostForm.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) AvailabilityJsonData(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}
	out, _ := json.MarshalIndent(resp, "", "   ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func (repo *Repository) MakeReservationPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "make-reservations.page.tmpl", r, &models.TemplateData{})
}
func (repo *Repository) ContactPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "contact.page.tmpl", r, &models.TemplateData{})
}
