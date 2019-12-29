package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID         string    `json:"id"`
	ResourceID string    `json:"resourceId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

func (r *Reservation) parseReservation(row *sql.Row) error {
	err := row.Scan(&r.ID, &r.ResourceID, &r.StartTime, &r.EndTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *Reservation) parseReservations(rows *sql.Rows) error {
	err := rows.Scan(&r.ID, &r.ResourceID, &r.StartTime, &r.EndTime)
	if err != nil {
		return err
	}
	return nil
}

// CreateTableReservations creates the table we use
func CreateTableReservations(db *sql.DB) (sql.Result, error) {
	q := `CREATE TABLE IF NOT EXISTS public.reservations
	(
		id uuid NOT NULL PRIMARY KEY,
		resource_id UUID REFERENCES resources(id),
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP NOT NULL
	)`

	return db.Exec(q)
}

func (r *Reservation) GetAllReservations(db *sql.DB) ([]*Reservation, error) {
	reservations := make([]*Reservation, 0)
	q := `SELECT * FROM reservations`

	rows, err := db.Query(q)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		r1 := new(Reservation)
		err := r1.parseReservations(rows)
		if err != nil {
			log.Fatal(err)
		}
		reservations = append(reservations, r1)
	}

	return reservations, nil
}

func (r *Reservation) GetReservationByID(db *sql.DB, id string) error {
	q := `SELECT * FROM reservations WHERE id = $1`

	row := db.QueryRow(q, id)
	err := r.parseReservation(row)
	if err != nil {
		return err
	}
	return nil
}

func (r *Reservation) GetReservationsByResourceID(db *sql.DB, id string) ([]*Reservation, error) {
	var err error
	var reservations []*Reservation
	q := `SELECT * FROM reservations WHERE resource_id = $1`

	rows, err := db.Query(q, id)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		r1 := new(Reservation)
		err := rows.Scan(&r1.ID, &r1.ResourceID, &r1.StartTime, &r1.EndTime)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, r1)
	}

	return reservations, nil
}

func (r *Reservation) CheckResourceAvailability(db *sql.DB) (bool, error) {
	var err error
	q := `SELECT * FROM reservations WHERE resource_id = $1 AND start_time >= $2 AND end_time <= $3`

	res, err := db.Exec(q, r.ResourceID, r.EndTime, r.StartTime)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if n > 0 {
		return false, nil
	}
	return true, nil
}

func (r *Reservation) NewReservation(db *sql.DB) error {
	var err error

	ok, err := r.CheckResourceAvailability(db)
	if err != nil {
		return err
	} else if !ok {
		return errors.New("room unavailable")
	}

	id, _ := uuid.NewRandom()
	q := `INSERT INTO reservations (id, resource_id, start_time, end_time) VALUES ($1, $2, $3, $4)`

	res, err := db.Exec(q, id, r.ResourceID, r.EndTime, r.StartTime)
	if err != nil {
		fmt.Println("error creating new reservation")
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		r.ID = id.String()
		return nil
	}
	return errors.New("error creating user")
}
