package db

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/xerrors"
	"github.com/jmoiron/sqlx"
)

// DBConn returns a database connection (or error if it can't connect) based on
// the given user, password and host. It will attempt to connect up to 20 times
// with an exponential back off.
func DBConn(user, password, dbName, host string) (*sqlx.DB, error) {
	maxRetries := 20
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, password, dbName, host)
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, xerrors.Errorf("failed to open database: %w", err)
	}
	for retry := 1; retry <= maxRetries; retry++ {
		log.Printf("Attempting to connect to database %s with user %s on host %s", dbName, user, host)
		if err = db.Ping(); err == nil {
			// If we can connect to the db okay then there is no point retrying
			// anymore so just exit here.
			return db, nil
		}
		log.Printf("Failed to connect to database on attempt %d: %v", retry, err)
		// Exponentially back off so we don't spam the db too much.
		time.Sleep(time.Duration(retry) * time.Second)
	}
	return nil, xerrors.Errorf("failed to connect to database %s after %d retries: %w", dbName, maxRetries, err)
}

// InitDB ensures the database is initialised with the required tables.
func InitDB(db *sqlx.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS name (
      id    SERIAL PRIMARY KEY
    , name  VARCHAR(246) NOT NULL CHECK (name <> '') UNIQUE
    , count INT DEFAULT 1
);

CREATE TABLE IF NOT EXISTS migration (
      id       SERIAL PRIMARY KEY
    , name     VARCHAR(246) NOT NULL CHECK (name <> '') UNIQUE
    , date_run TIMESTAMP DEFAULT NOW()
);
`
	_, err := db.Exec(schema)
	if err != nil {
		return xerrors.Errorf("failed to create database tables from schema: %w", err)
	}

	migrations := []Migration{}
	ApplyMigration(db, migrations...)
	if err != nil {
		return xerrors.Errorf("failed to migrate count default from 0 to 1: %w", err)
	}
	return nil
}
