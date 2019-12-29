package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// Resource is a model for the resources
type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateTableResource(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS public.resources
	(
		id UUID NOT NULL PRIMARY KEY,
		name VARCHAR(30) UNIQUE NOT NULL
	)`

	_, err := db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resource) CreateResource(db *sql.DB) error {
	var err error
	id, _ := uuid.NewRandom()
	q := `INSERT INTO resources (id, name) VALUES ($1, $2)`

	_, err = db.Exec(q, id, r.Name)
	if err != nil {
		return err
	}
	r.ID = id.String()
	return nil
}

func (r *Resource) DeleteResource(db *sql.DB) error {
	return errors.New("'delete records' not implemented")
}

func (r *Resource) GetAllResources(db *sql.DB) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	q := `SELECT * FROM resources`

	rows, err := db.Query(q)
	if err != nil {
		return resources, err
	}

	for rows.Next() {
		resource := new(Resource)
		rows.Scan(&resource.ID, &resource.Name)
		resources = append(resources, resource)
	}
	return resources, nil
}

func (r *Resource) GetResourceByID(db *sql.DB, id string) error {
	q := `SELECT * FROM resources WHERE id = $1`

	row := db.QueryRow(q, id)
	err := row.Scan(&r.ID, &r.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resource) FindResource(db *sql.DB, query string) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	q := `SELECT * FROM resources WHERE name LIKE $1`

	prepped := query + "%"

	rows, err := db.Query(q, prepped)
	if err != nil {
		return resources, err
	}
	for rows.Next() {
		var r1 = new(Resource)
		err := rows.Scan(&r1.ID, &r1.Name)
		if err != nil {
			return resources, err
		}
		resources = append(resources, r1)
	}
	return resources, nil
}
