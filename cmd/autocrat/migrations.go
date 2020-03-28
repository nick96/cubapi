package main

import (
	"time"

	"github.com/nick96/cubapi/db/migrate"
)

const (
	// Latest version of the database schema. This will be applied to new
	// database so the migrations will not be required.
	Schema = `
CREATE TABLE IF NOT EXISTS users (

);`
)

var (
	// Migrations to applied to an existing database.
	Migrations = []migrate.Migration{
		{
			Version: 1,
			Date: time.Date(2020, 03, 22, 22, 59, 0, 0, time.FixedZone("Australia/Melbourne", 10 * time.Hour)),
			SQL: `
CREATE TABLE IF NOT EXISTS users (

);
`,
			Description: "Initialise required tables.",
		}
	}
)
