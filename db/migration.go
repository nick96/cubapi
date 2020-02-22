package main

import (
	"crypto/sha256"
	"database/sql/driver"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

// SHA256Checksum represents a 256 byte checksum created using SHA256.
type SHA256Checksum [sha256.Size]byte

func (c SHA256Checksum) Scan(src interface{}) error {
	return nil
}

func (c SHA256Checksum) Value() (driver.Value, error) {
	return fmt.Sprintf("%x", c), nil
}

// Migration represents a migration.
type Migration struct {
	// Name is a unique name to identify the migration (both for the human and
	// the machine).
	Name string
	// Date is the date and time the migration was written.
	Date time.Time
	// Script is the migration script to be applied.
	Script string
}

// MigrationDB represents the migration entity in the database.
type MigrationDB struct {
	// Name is a unique name to identify the migration (both for the human and
	// the machine).
	Name string
	// Date is the date and time the migration was written.
	Date time.Time
	// DateApplied is the date the migration was applied to the database.
	DateApplied time.Time
	// Checksum is a checksum of the migration script to ensure it has not been
	// changed.
	Checksum SHA256Checksum
}

type Migrations []Migration

func InitMigration(db *sqlx.DB) error {
	migrationTable := `
CREATE TABLE IF NOT EXISTS migration (
    id           SERIAL       PRIMARY KEY
  , name         VARCHAR(246) UNIQUE NOT NULL
  , date         TIMESTAMP    NOT NULL
  , date_applied TIMESTAMP    NOT NULL
  , checksum     VARCHAR(256) NOT NULL
);
`
	if _, err := db.Exec(migrationTable); err != nil {
		return fmt.Errorf("Failed to initialise databse with migration table: %w", err)
	}
	return nil
}

// Apply applies a given migration to the supplied transaction handle.
func (m Migration) Apply(tx *sqlx.Tx) error {
	if applied, err := m.isApplied(tx); err != nil {
		return fmt.Errorf("Failed to check if migration '%s' was applied: %v", m.Name, err)
	} else if !applied {
		return m.apply(tx)
	}
	log.Printf("Migration '%s' was not applied as it has already been applied on this database", m.Name)
	return nil
}

// isApplied checks if a migration has been applied.
func (m Migration) isApplied(tx *sqlx.Tx) (bool, error) {
	var exists bool
	err := tx.Get(&exists, `SELECT EXISTS (SELECT * FROM migration WHERE name = $1);`, m.Name)
	if err != nil {
		return false, fmt.Errorf("failed to check if migration '%s' as been applied: %w", m.Name, err)
	}
	return exists, nil
}

// apply applys a database migration and marks it as completed if successful.
func (m Migration) apply(tx *sqlx.Tx) error {
	log.Printf("Applying migration '%s' added on %v", m.Name, m.Date)
	_, err := tx.Exec(m.Script)
	if err != nil {
		return fmt.Errorf("failed to apply migration '%s': %v", m.Name, err)
	}

	err = m.markCompleted(tx)
	if err != nil {
		return fmt.Errorf("failed to mark migration '%s' as completed: %v", m.Name, err)
	}
	return nil
}

// markComplete inserts a migration into the database to indicate it is
// completed.
func (m Migration) markCompleted(tx *sqlx.Tx) error {
	checksum := sha256.Sum256([]byte(m.Script))
	migrationDB := MigrationDB{
		Name:        m.Name,
		Date:        m.Date,
		DateApplied: time.Now(),
		Checksum:    SHA256Checksum(checksum),
	}
	_, err := tx.Exec(`INSERT INTO migration(name, date, date_applied, checksum) VALUES ($1, $2, $3, $4)`,
		migrationDB.Name, migrationDB.Date, migrationDB.DateApplied, migrationDB.Checksum)
	return err
}

func (m Migrations) Len() int {
	return len(m)
}

func (m Migrations) Less(i, j int) bool {
	return m[i].Date.Before(m[j].Date)
}

func (m Migrations) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}

// ApplyMigration applies the given migrations to the database. Migrations are
// applied in order from the oldest migration to the newest.
func ApplyMigration(db *sqlx.DB, migrations ...Migration) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to initialise transation: %w", err)
	}

	sort.Sort(Migrations(migrations))
	for _, migration := range migrations {
		if err = migration.Apply(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to apply migration '%s': %w", migration.Name, err)
		}
	}

	tx.Commit()
	return nil
}
