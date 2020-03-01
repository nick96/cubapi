package attendance

import (
	"github.com/jmoiron/sqlx"
)

// Cub represents a cub.
type Cub struct {
	Model

	// FirstName is the first name of the cub.
	FirstName string `json:"first_name"`
	// LastName is the last name of the cub.
	LastName string `json:"last_name"`
	// Attendances is the list of zero or more recorded attendances for the cub.
	Attendances []Attendance `json:"attendances"`
}

type CubStoreReader interface{}

type CubStore struct {
	db *sqlx.DB
}

func NewCubStore(db *sqlx.DB) CubStore {
	return CubStore{db}
}

type CubsHandler struct {
	cubStore CubStoreReader
}

