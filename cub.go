package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/nick96/cubapi/model"
)

// Cub represents a cub.
type Cub struct {
	model.Model

	// FirstName is the first name of the cub.
	FirstName string `json:"first_name"`
	// LastName is the last name of the cub.
	LastName string `json:"last_name"`
	// Attendances is the list of zero or more recorded attendances for the cub.
	Attendances []Attendance `json:"attendances"`
}

type CubStore struct {
	DB *sqlx.DB
}

type CubsHandler struct {
	cubStore CubStoreReader
}

type CubHandler struct {
	cubStore CubStoreReader
}


