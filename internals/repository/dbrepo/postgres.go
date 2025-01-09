package dbrepo

import (
	"context"
	"github.com/nicholas-karimi/bookings/internals/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation insers a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var newID int

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date,
                          room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	//_, err := m.DB.ExecContext(ctx, stmt,
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (m *postgresDBRepo) RoomExists(roomID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM rooms WHERE id = $1)`
	err := m.DB.QueryRow(query, roomID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// InsertRoomRestriction() insert a room restriction into db
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}
