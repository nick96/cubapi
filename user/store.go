package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotImplemented = errors.New("Not Implemented")
)

// UserStorer is an interface that must be implemented by things that store user
// information.
type UserStorer interface {
	FindByEmail(email string) (User, bool, error)
	AddUser(user User) (int64, error)
}

// UserStore is a store for users and their related information. It implements
// the UserStorer interface.
type UserStore struct {
	db *sqlx.DB
}

// NewStore creates a new store from the given sqlx db handle.
func NewStore(db *sqlx.DB) UserStorer {
	return UserStore{db}
}

// FindByEmail finds a user by their email.
func (s UserStore) FindByEmail(email string) (user User, found bool, err error) {
	query := `SELECT * FROM autocrat.users WHERE email = $1;`
	err = s.db.QueryRowx(query, email).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false, nil
		}
		return User{}, false, fmt.Errorf("could not find user with email '%s': %w", email, err)
	}
	return user, true, nil
}

// AddUser adds the given user to the database and returns the ID of the
// inserted user.
func (s UserStore) AddUser(user User) (int64, error) {
	var id int64
	query := `
	INSERT INTO autocrat.users (id, email, firstname, lastname, password)
	VALUES (DEFAULT, $1, $2, $3, $4) 
	RETURNING id;
	`
	err := s.db.
		QueryRow(query, user.Email, user.FirstName, user.LastName, user.Password).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user into store: %w", err)
	}
	return id, nil
}
