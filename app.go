package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	// dbHost     = "DBHOST"
	// dbPort     = "DBPORT"
	dbUser     = "DB_USER"
	dbName     = "DB_NAME"
	dbPassword = "DB_PASSWORD"
)

type app struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *app) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to DB!")

	a.Router = mux.NewRouter()

	fmt.Println("Tables created...")

	api := a.Router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server healthy"))
	}).Methods(http.MethodGet)

	api.HandleFunc("/reservations/resource", a.getResource).Methods(http.MethodGet)
	api.HandleFunc("/reservations/resource/new", a.createResource).Methods(http.MethodPost)

	api.HandleFunc("/reservations/new", a.newReservation).Methods(http.MethodPost)
	api.HandleFunc("/reservations", a.getAllReservations).Methods(http.MethodGet)

	//	TODO: 	Check availablity

	a.Router.Use(middlewareLogger)
}

func (a *app) Run(addr string) {
	fmt.Printf("Server listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
