package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nicholas-karimi/bookings/internals/driver"
	"github.com/nicholas-karimi/bookings/internals/helpers"
	"github.com/nicholas-karimi/bookings/internals/repository"
	"github.com/nicholas-karimi/bookings/internals/repository/dbrepo"
	"net/http"
	"strconv"
	"time"

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
	DB  repository.DatabaseRepo
}

// cretaes a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets repo forthe handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {

	//repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, "home.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "about.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) IndexPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "index.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) GeneralsPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) MajorsPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) AvailabilityPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.tmpl", r, &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.PostForm.Get("start")
	end := r.PostForm.Get("end")

	// 2024-01-01 -> cast to this - 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetch rooms availability
	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		fmt.Println("Error fetching room availability:", err)
		return
	}

	//fmt.Println("Found rooms:", len(rooms))
	//for _, room := range rooms {
	//	//fmt.Println("Room:", room) // Print room details for debugging
	//	repo.App.InfoLog.Println("ROOM", room.ID, room.RoomName)
	//}
	if len(rooms) == 0 {
		// no availability
		//repo.App.InfoLog.Println("No rooms available")
		repo.App.Session.Put(r.Context(), "error", "No rooms available for dates selected.")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms
	//w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	repo.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, "choose_room.page.tmpl", r, &models.TemplateData{
		Data: data,
	})
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

	render.Template(w, "make-reservations.page.tmpl", r, &models.TemplateData{

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

	// Parse Dates
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2024-01-01 -> cast to this - 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//Parse room id
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Check if the room exists
	//exists, err := dbrepo.RoomExists(roomID)
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	//if !exists {
	//	http.Error(w, "Invalid room ID", http.StatusBadRequest)
	//	return
	//}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}
	// Validate the form
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
		render.Template(w, "make-reservations.page.tmpl", r, &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Insert reservation into DB
	newReservationID, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{

		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}
	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Store reservation in session and redirect

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
	render.Template(w, "reservation-summary.page.tmpl", r, &models.TemplateData{
		Data: data,
	})
}
func (repo *Repository) ContactPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl", r, &models.TemplateData{})
}
