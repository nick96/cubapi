package user

import (
	"github.com/jmoiron/sqlx"
)

type UserStorer interface{}

type UserStore struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) UserStorer {
	return UserStore{db}
}
