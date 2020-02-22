package main

import (
	"flag"
	"fmt"
	"log"
	"testing"

	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

var dbHandle *sqlx.DB

func TestMain(m *testing.M) {
	// These tests take a long time so we don't want to run them if `-short` has
	// been used.
	flag.Parse()
	if !testing.Short() {
		pool, err := dockertest.NewPool("")
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
			dbHandle, err = sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/db?sslmode=disable", resource.GetPort("5432/tcp")))
			if err != nil {
				return err
			}
			return dbHandle.Ping()
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer pool.Purge(resource)

		cleanup := setup()
		defer cleanup()
		m.Run()
	}
}

// setup sets up the environment for testing.
func setup() func() {
	schema := `
CREATE TABLE IF NOT EXISTS table1 (id SERIAL PRIMARY KEY, x INT DEFAULT 0);
CREATE TABLE IF NOT EXISTS table2 (id SERIAL PRIMARY KEY, y INT DEFAULT 1);
`
	dbHandle.MustExec(schema)

	InitMigration(dbHandle)

	return func() {
		stmts := `
DROP TABLE IF EXISTS table1;
DROP TABLE IF EXISTS table2;
DROP TABLE IF EXISTS migration;
`
		dbHandle.MustExec(stmts)
	}
}

func TestInitMigrationNonExistantTable(t *testing.T) {
	// Ensure the table doesn't exist
	dbHandle.MustExec(`DROP TABLE IF EXISTS migration;`)
	var exists bool
	dbHandle.Get(&exists, `SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = 'migration');`)
	if exists {
		t.Fatalf("Expected 'migration' table not to exist after drop but it does")
	}

	// Initialise and check there is no problem
	err := InitMigration(dbHandle)
	if err != nil {
		t.Fatalf("Expected migration table to be initialised successfully: %v", err)
	}

	// Check the table now exists
	dbHandle.Get(&exists, `SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = 'migration');`)
	if !exists {
		t.Fatalf("Expected 'migration' table to exist but it doesn't")
	}
}

func TestInitMigrationExistingTable(t *testing.T) {
	cleanup := setup()
	defer cleanup()

	// Ensure the table exists
	err := InitMigration(dbHandle)
	if err != nil {
		t.Fatalf("Failed to initialise migration table: %v", err)
	}

	// Initialise and check there is no problem
	err = InitMigration(dbHandle)
	if err != nil {
		t.Fatalf("Expected migration table to be initialised successfully: %v", err)
	}

	// Check the table now exists
	var exists bool
	dbHandle.Get(&exists, `SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = 'migration');`)
	if !exists {
		t.Fatalf("Expected 'migration' table to exist but it doesn't")
	}
}

func TestApplyMigration(t *testing.T) {
	cleanup := setup()
	defer cleanup()

	var beforeID int
	err := dbHandle.QueryRow(`INSERT INTO table1 VALUES (DEFAULT, DEFAULT) RETURNING id;`).Scan(&beforeID)
	if err != nil {
		t.Fatalf("Expected to retrieve before ID without error: %v", err)
	}

	migration := Migration{
		Name:   "migration-1",
		Date:   time.Date(2009, 11, 17, 20, 34, 58, 12134324, time.UTC),
		Script: `ALTER TABLE table1 ALTER COLUMN x SET DEFAULT 1;`,
	}
	tx := dbHandle.MustBegin()
	err = migration.Apply(tx)
	if err != nil {
		t.Fatalf("Expected migration to be applied without error: %v", err)
	}
	tx.Commit()

	var afterID int
	err = dbHandle.QueryRow(`INSERT INTO table1 VALUES (DEFAULT, DEFAULT) RETURNING id;`).Scan(&afterID)
	if err != nil {
		t.Fatalf("Expected to retrieve after ID without error: %v", err)
	}

	var defaultBefore, defaultAfter int
	dbHandle.Get(&defaultBefore, "SELECT x FROM table1 WHERE id = $1", beforeID)
	dbHandle.Get(&defaultAfter, "SELECT x FROM table1 WHERE id = $1", afterID)

	if defaultBefore != 0 {
		t.Errorf("Expected default value for 'x' before migration to be 0, found %d", defaultBefore)
	}

	if defaultAfter != 1 {
		t.Errorf("Expected default value for 'x' after migration to be 1, found %d", defaultAfter)
	}
}

func TestApplyMultipleMigrations(t *testing.T) {
	cleanup := setup()
	defer cleanup()

	var beforeIDTable1 int
	err := dbHandle.QueryRow(`INSERT INTO table1 VALUES (DEFAULT, DEFAULT) RETURNING id;`).Scan(&beforeIDTable1)
	if err != nil {
		t.Fatalf("Expected to retrieve before ID table 1 without error: %v", err)
	}

	var beforeIDTable2 int
	err = dbHandle.QueryRow(`INSERT INTO table2 VALUES (DEFAULT, DEFAULT) RETURNING id;`).Scan(&beforeIDTable2)
	if err != nil {
		t.Fatalf("Expected to retrieve before ID for table 2 without error: %v", err)
	}

	migrations := []Migration{
		Migration{
			Name:   "migration-1",
			Date:   time.Date(2009, 11, 17, 20, 34, 58, 12134324, time.UTC),
			Script: `ALTER TABLE table1 ALTER COLUMN x SET DEFAULT 1;`,
		},
		Migration{
			Name:   "migration-2",
			Date:   time.Date(2009, 11, 18, 20, 34, 58, 12134324, time.UTC),
			Script: `ALTER TABLE table2 ALTER COLUMN y SET DEFAULT 2;`,
		},
	}
	err = ApplyMigration(dbHandle, migrations...)
	if err != nil {
		t.Fatalf("Expected migration to be applied without error: %v", err)
	}

	var afterIDTable1 int
	err = dbHandle.QueryRow(`INSERT INTO table1 VALUES (DEFAULT, DEFAULT) RETURNING id;`).Scan(&afterIDTable1)
	if err != nil {
		t.Fatalf("Expected to retrieve after ID table 2 without error: %v", err)
	}

	var afterIDTable2 int
	err = dbHandle.QueryRow(`INSERT INTO table2 VALUES (DEFAULT, DEFAULT) RETURNING id;`).Scan(&afterIDTable2)
	if err != nil {
		t.Fatalf("Expected to retrieve after ID for table 2 without error: %v", err)
	}

	var defaultBeforeTable1, defaultAfterTable1 int
	err = dbHandle.Get(&defaultBeforeTable1, "SELECT x FROM table1 WHERE id = $1", beforeIDTable1)
	if err != nil {
		t.Fatalf("Expected to retrieve x value from table 1 without error: %v", err)
	}
	err = dbHandle.Get(&defaultAfterTable1, "SELECT x FROM table1 WHERE id = $1", afterIDTable1)
	if err != nil {
		t.Fatalf("Expected to retrieve x value from table 1 without error: %v", err)
	}

	var defaultBeforeTable2, defaultAfterTable2 int
	err = dbHandle.Get(&defaultBeforeTable2, "SELECT y FROM table2 WHERE id = $1", beforeIDTable2)
	if err != nil {
		t.Fatalf("Expected to retrieve y value from table 2 without error: %v", err)
	}
	err = dbHandle.Get(&defaultAfterTable2, "SELECT y FROM table2 WHERE id = $1", afterIDTable2)
	if err != nil {
		t.Fatalf("Expected to retrieve y value from table 2 without error: %v", err)
	}

	if defaultBeforeTable1 != 0 {
		t.Errorf("Expected default value for 'x' before migration to be 0, found %d", defaultBeforeTable1)
	}

	if defaultAfterTable1 != 1 {
		t.Errorf("Expected default value for 'x' after migration to be 1, found %d", defaultAfterTable1)
	}

	if defaultBeforeTable2 != 1 {
		t.Errorf("Expected default value for 'y' before migration to be 1, found %d", defaultBeforeTable2)
	}

	if defaultAfterTable2 != 2 {
		t.Errorf("Expected default value for 'y' after migration to be 2, found %d", defaultAfterTable2)
	}
}

func TestDoNotApplyAppliedMigration(t *testing.T) {
	cleanup := setup()
	defer cleanup()

	migration := Migration{
		Name:   "migration-1",
		Date:   time.Date(2009, 11, 17, 20, 34, 58, 12134324, time.UTC),
		Script: `INSERT INTO table1 VALUES (DEFAULT, DEFAULT);`,
	}
	tx := dbHandle.MustBegin()
	err := migration.Apply(tx)
	if err != nil {
		t.Fatalf("Expected migration to be applied without error: %v", err)
	}
	tx.Commit()

	tx = dbHandle.MustBegin()
	err = migration.Apply(tx)
	if err != nil {
		t.Fatalf("Expected migration to be applied without error: %v", err)
	}
	tx.Commit()

	var count int
	err = dbHandle.QueryRow(`SELECT COUNT(*) FROM table1;`).Scan(&count)
	if count != 1 {
		t.Fatalf("Expected migration to only be applied once")
	}
}

func TestReturnMigrationError(t *testing.T) {
	cleanup := setup()
	defer cleanup()
}

func TestRollbackAllMigrationsOnError(t *testing.T) {
	cleanup := setup()
	defer cleanup()
}
