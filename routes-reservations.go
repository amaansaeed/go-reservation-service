package main

import (
	"encoding/json"
	"net/http"
	"github.com/amaansaeed/go-reservation-service/models"
	"time"
)

type newReservation struct {
	ResourceID string `json:"resourceId"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
}

func (a *app) newReservation(w http.ResponseWriter, r *http.Request) {
	var nr newReservation
	var reservation models.Reservation
	var err error

	err = json.NewDecoder(r.Body).Decode(&nr)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Issue with incoming JSON\n"))
		return
	}
	if nr.ResourceID == "" || nr.StartTime == 0 || nr.EndTime == 0 || nr.StartTime > nr.EndTime {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("improper dataset\n"))
		return
	}

	reservation.ResourceID = nr.ResourceID
	reservation.StartTime = time.Unix(nr.StartTime, 0)
	reservation.EndTime = time.Unix(nr.EndTime, 0)

	err = reservation.NewReservation(a.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	res, _ := json.Marshal(reservation)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *app) getAllReservations(w http.ResponseWriter, r *http.Request) {
	var r1 models.Reservation
	var err error

	reservations, err := r1.GetAllReservations(a.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	res, _ := json.Marshal(reservations)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}
