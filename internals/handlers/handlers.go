package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nicholas-karimi/bookings/internals/helpers"
	"net/http"

	"github.com/nicholas-karimi/bookings/internals/config"
	"github.com/nicholas-karimi/bookings/internals/forms"
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

	//repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplates(w, "home.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplates(w, "about.page.tmpl", r, &models.TemplateData{})
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
	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func (repo *Repository) MakeReservationPage(w http.ResponseWriter, r *http.Request) {
	// renepopulate the form with data when theres an error
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplates(w, "make-reservations.page.tmpl", r, &models.TemplateData{

		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservationPage will handle the posting of the form
func (repo *Repository) PostReservationPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// could not parse form
		//w.Write([]byte("ParseForm() err: " + err.Error()))
		helpers.ServerError(w, err)
		return

	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// create a new form
	form := forms.New(r.PostForm)

	// check if form is valid
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email", r)
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplates(w, "make-reservations.page.tmpl", r, &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// pass the date to a sessiom

	repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		//log.Println("cannot get item from reservation")
		repo.App.ErrorLog.Println("Cannot get error from session")
		repo.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	// remove the reservation from the session
	repo.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplates(w, "reservation-summary.page.tmpl", r, &models.TemplateData{
		Data: data,
	})
}
func (repo *Repository) ContactPage(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "contact.page.tmpl", r, &models.TemplateData{})
}
