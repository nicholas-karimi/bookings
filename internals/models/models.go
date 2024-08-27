package models


// Reservation hold data sent from handlers to templates
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	// StartDate time.Time
	// EndDate   time.Time
	// RoomID    string
	// Room      Room
}