package repository

import "github.com/nicholas-karimi/bookings/internals/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(reservation models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
