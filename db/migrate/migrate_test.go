package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"go.uber.org/zap"
)

var db *sql.DB

func checkTableExists(tableName string) (bool, error) {
	var exists bool
	row := db.QueryRow(`SELECT EXISTS (
SELECT FROM information_schema.tables
WHERE table_name = $1
)`, tableName)
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check if table '%s' exists: %v", tableName, err)
	}
	return exists, nil
}

func cleanup() {
	_, err := db.Exec(`DROP SCHEMA public CASCADE;`)
	if err != nil {
		log.Fatalf("Failed to drop public schema: %v", err)
	}
	_, err = db.Exec(`CREATE SCHEMA public;`)
	if err != nil {
		log.Fatalf("Failed to recreate public schema: %v", err)
	}
}

func TestMain(m *testing.M) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed to create dockertest pool: %v", err)
	}

	pgUser := "test"
	pgPass := "password"
	pgDB := "testdb"
	envVars := []string{
		fmt.Sprintf("POSTGRES_USER=%s", pgUser),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPass),
		fmt.Sprintf("POSTGRES_DB=%s", pgDB),
	}
	resource, err := pool.Run("postgres", "", envVars)
	if err != nil {
		log.Fatalf("Failed to start dockertest resource: %v", err)
	}

	connString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgUser, pgPass, pgDB, resource.GetPort("5432/tcp"),
	)
	err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", connString)
		if err != nil {
			return fmt.Errorf("failed to open connection to database %s (%s): %v", pgDB, connString, err)
		}

		if err = db.Ping(); err != nil {
			return fmt.Errorf("failed to ping database %s (%s): %v", pgDB, connString, err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Could not connect to docker db: %v", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}

	os.Exit(code)
}

func TestInit(t *testing.T) {
	logger := zap.NewNop()
	migrator := NewMigrator(db, logger)
	exists, err := checkTableExists("migrations")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("Expected 'migrations' table not to exist but it does")
	}

	if err := migrator.Init(); err != nil {
		t.Fatalf("Failed to apply migration Init to db: %v", err)
	}

	exists, err = checkTableExists("migrations")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("Expected 'migrations' table to exist but it does not")
	}
}

func TestMigrate(t *testing.T) {
	tests := []struct {
		Name       string
		Migrations []Migration
	}{
		{
			Name: "single-migration",
			Migrations: []Migration{
				{
					Version: 1,
				},
			},
		},
		{
			Name: "multi-migration-all-applied",
		},
		{
			Name: "multi-migration-some-applied",
		},
	}

	logger := zap.NewNop()
	migrator := NewMigrator(db, logger)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Cleanup(cleanup)
			err := migrator.Apply(tt.Migrations...)
			if err != nil {
				t.Fatalf("Expected migration application to succeed: %v", err)
			}
		})
	}
}

func TestMigrateFail(t *testing.T) {
	tests := []struct {
		Name       string
		Migrations []Migration
		Check      func() error
	}{
		{
			Name: "single-migration-failure",
			Migrations: []Migration{
				{
					Version:     1,
					Date:        time.Now(),
					SQL:         `CREATE TABLE IF NOT EXITS test;`,
					Description: "Migration that fails due to a syntax error",
				},
			},
			Check: func() error { return nil },
		},
		{
			Name: "multi-migration-all-failure",
			Migrations: []Migration{
				{
					Version:     1,
					Date:        time.Now(),
					SQL:         `CREATE TABLE IF NOT EXITS test;`,
					Description: "Migration that fails due to a syntax error",
				},
				{
					Version:     2,
					Date:        time.Now(),
					SQL:         `CREATE TABLE IF NOT EXITS test;`,
					Description: "Migration that fails due to a syntax error",
				},
			},
			Check: func() error { return nil },
		},
		{
			Name: "multi-migration-single-failure",
			Migrations: []Migration{
				{
					Version:     1,
					Date:        time.Now(),
					SQL:         `CREATE TABLE IF NOT EXISTS test;`,
					Description: "Migration that creates a table",
				},
				{
					Version:     2,
					Date:        time.Now(),
					SQL:         `CREATE TABLE IF NOT EXITS test;`,
					Description: "Migration that fails due to a syntax error",
				},
			},
			Check: func() error {
				exists, err := checkTableExists("test")
				if err != nil {
					return err
				}
				if exists {
					return fmt.Errorf("Expected 'test' table not to exist but it does")
				}
				return nil
			},
		},
	}

	logger := zap.NewNop()
	migrator := NewMigrator(db, logger)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := migrator.Apply(tt.Migrations...)
			if err == nil {
				t.Fatalf("Expected migration application to fail")
			}
			err = tt.Check()
			if err != nil {
				t.Fatalf("Expected check to pass: %v", err)
			}
		})
	}

}

func TestSchema(t *testing.T) {
	logger := zap.NewNop()
	migrator := NewMigrator(db, logger)

	exists, err := checkTableExists("test")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("Expected 'test' table not to exist but it does")
	}

	if err := migrator.Schema(`CREATE TABLE test (id INT);`); err != nil {
		t.Fatal(err)
	}

	exists, err = checkTableExists("test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("Expected 'test' table to exists but it does not")
	}

}
