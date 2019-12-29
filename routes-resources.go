package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"github.com/amaansaeed/go-reservation-service/models"
	"github.com/amaansaeed/go-reservation-service/utils"
)

type newResource struct {
	Name string `json:"name"`
}

func (a *app) getResource(w http.ResponseWriter, r *http.Request) {
	var resource models.Resource

	id := r.URL.Query().Get("id")
	q := r.URL.Query().Get("q")

	if (id == "" || !utils.IsValidUUID(id)) && q == "" {
		resources, err := resource.GetAllResources(a.DB)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error getting all resources"))
			return
		}
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(resources)
		w.Write(res)
		return
	}

	if len(id) > 1 {
		err := resource.GetResourceByID(a.DB, id)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no resource found"))
			return
		} else if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(resource)
		w.Write(res)
		return

	} else if len(q) > 1 {
		resources, err := resource.FindResource(a.DB, q)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no resource found"))
			return
		} else if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(resources)
		w.Write(res)
		return
	}
	return
}

func (a *app) createResource(w http.ResponseWriter, r *http.Request) {
	var nr newResource
	var resource models.Resource
	var err error

	err = json.NewDecoder(r.Body).Decode(&nr)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Issue with incoming JSON"))
		return
	}
	if nr.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name missing\n"))
		return
	}

	resource.Name = nr.Name

	err = resource.CreateResource(a.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not create resource"))
		return
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(resource)
	w.Write(res)
}
