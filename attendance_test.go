package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/nick96/cubapi/db"
	"github.com/org/dockertest"
)

var dbHandle *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("attendance")
	if err != nil {
		log.Fatalf("Failed to create docker pool: %v", err)
	}

	resource, err := pool.Run("postgres", "9.6.17", []string{
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=db",
	})
	if err != nil {
		log.Fatalf("Failed to start resource: %v", err)
	}

	err = pool.Retry(func() error {
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/db", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Purge(resource)

	err = db.InitDB(dbHandle)
	if err != nil {
		log.Fatalf("Failed to initialise the database: %v", err)
	}

	m.Run()
}

func TestGetAll(t *testing.T) {
}

func TestGetByID(t *testing.T) {

}

func TestGetByCub(t *testing.T) {

}

func TestGetByDate(t *testing.T) {
}
