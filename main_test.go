package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var router = mux.NewRouter()

var a app

func TestMain(m *testing.M) {
	godotenv.Load()
	a = app{}
	a.Initialize(os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		fmt.Sprintf("%s_test", os.Getenv("DB_NAME")))
	createTables()
	seedTestData()
	code := m.Run()
	destroyTable()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func prepPayload(p string) io.Reader {
	return bytes.NewBuffer([]byte(p))
}

func createTables() {
	q1 := `CREATE TABLE IF NOT EXISTS public.resources
	(
		id uuid NOT NULL PRIMARY KEY,
		name VARCHAR(30) UNIQUE NOT NULL
	)`
	q2 := `CREATE TABLE IF NOT EXISTS public.reservations
	(
		id uuid NOT NULL PRIMARY KEY,
		resource_id uuid REFERENCES resources(id),
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP NOT NULL
	)`

	a.DB.Exec(q1)
	a.DB.Exec(q2)
}

func seedTestData() {
	q1 := `INSERT INTO resources (id, name) VALUES ($1, $2)`
	q2 := `INSERT INTO reservations (id, resource_id, start_time, end_time) VALUES ($1, $2, $3, $4)`

	id1 := "c3bd1ddc-fa07-4edc-b268-15cc18c95f01"
	id2 := "eb2fce55-a6e4-4275-9115-7e5bb918df9e"
	id3, _ := uuid.NewRandom()
	id4, _ := uuid.NewRandom()
	id5, _ := uuid.NewRandom()

	st1 := time.Unix(1577293360, 0)
	et1 := time.Unix(1578293360, 0)
	st2 := time.Unix(1597293360, 0)
	et2 := time.Unix(1597293360, 0)

	a.DB.Exec(q1, id1, "room 1")
	a.DB.Exec(q1, id2, "room 2")
	a.DB.Exec(q1, id3, "room 3")

	a.DB.Exec(q2, id4, id1, st1, et1)
	a.DB.Exec(q2, id5, id2, st2, et2)
}

func destroyTable() {
	a.DB.Exec("DROP TABLE IF EXISTS public.resources")
	a.DB.Exec("DROP TABLE IF EXISTS public.reservations")
}

func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "server healthy" {
		t.Errorf("Expected 'server healthy'. Got %s", body)
	}
}

func TestResources(t *testing.T) {
	id1 := "c3bd1ddc-fa07-4edc-b268-15cc18c95f01"

	req1, _ := http.NewRequest("GET", "/api/v1/reservations/resource", nil)
	res1 := executeRequest(req1)

	checkResponseCode(t, http.StatusOK, res1.Code)

	req2, _ := http.NewRequest("GET", "/api/v1/reservations/resource?id="+id1, nil)
	res2 := executeRequest(req2)

	checkResponseCode(t, http.StatusOK, res2.Code)

	req3, _ := http.NewRequest("GET", "/api/v1/reservations/resource?id=invalid", nil)
	res3 := executeRequest(req3)

	checkResponseCode(t, http.StatusOK, res3.Code)

	req4, _ := http.NewRequest("GET", "/api/v1/reservations/resource?q=room+1", nil)
	res4 := executeRequest(req4)

	checkResponseCode(t, http.StatusOK, res4.Code)

	// if body := response.Body.String(); body != "server healthy" {
	// 	t.Errorf("Expected 'server healthy'. Got %s", body)
	// }
}

func TestCreateNewResource(t *testing.T) {
	p := `{"name": "book 1"}`

	req1, _ := http.NewRequest("POST", "/api/v1/reservations/resource/new", prepPayload(p))
	res1 := executeRequest(req1)

	checkResponseCode(t, http.StatusOK, res1.Code)
}

func TestGetReservations(t *testing.T) {
	req1, _ := http.NewRequest("GET", "/api/v1/reservations", nil)
	res1 := executeRequest(req1)

	checkResponseCode(t, http.StatusOK, res1.Code)

	// if body := response.Body.String(); body != "server healthy" {
	// 	t.Errorf("Expected 'server healthy'. Got %s", body)
	// }
}
