package main

import (
	"time"

	"github.com/nick96/cubapi/db/migrate"
)

var (
	// Migrations to applied to an existing database.
	Migrations = []migrate.Migration{
		{
			Version: 1,
			Date:    time.Date(2020, 04, 12, 12, 23, 0, 0, time.FixedZone("Australia/Melbourne", 10)),
			SQL: `
CREATE SCHEMA autocrat
    CREATE TABLE users (
          id           SERIAL       PRIMARY KEY
        , email        VARCHAR(256) NOT NULL
        , firstname    VARCHAR(256) NOT NULL
        , lastname     VARCHAR(256) NOT NULL
        , password     CHAR(60)     NOT NULL
        , salt         CHAR(10)     NOT NULL
    );
`,
			Description: "Initial schema.",
		},
		{
			Version: 2,
			Date:    time.Date(2020, 04, 12, 14, 38, 0, 0, time.FixedZone("Australia/Melbourne", 10)),
			SQL: `
ALTER TABLE autocrat.users
    DROP COLUMN salt;
`,
			Description: "Remove salt column as we're using bcrypt which generates the salt as part of the hash.",
		},
	}
)
