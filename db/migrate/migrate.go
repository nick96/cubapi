package migrate

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// Migrator is an applyer of migrations.
type Migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

// Migration is a migration that should be applied to the database.
type Migration struct {
	// Version is the migration version. This should be unique across all migrations.
	Version int
	// Date the migration was written.
	Date time.Time
	// SQL script to apply as part of the migration.
	SQL string
	// Description of the migration. This is not required but can sometimes
	// be useful to give context to a complex migration.
	Description string
}

// AppliedMigration is a migration that has been applied. It represents a row in
// the schema history table in the database.
type AppliedMigration struct {
	// Version is the migration version.
	Version int
	// DateCreated is the date a migration was created.
	DateCreated time.Time
	// DateApplied is the data a migration was applied to the database.
	DateApplied time.Time
	// Description is a description of the migration.
	Description string
	// Checksum is a SHA256 checksum of the SQL script applied in the
	// migration.
	Checksum []byte
}

// NewMigrator returns a new migrator with the given DB handle and logger.
func NewMigrator(db *sql.DB, logger *zap.Logger) Migrator {
	return Migrator{
		db:     db,
		logger: logger,
	}
}

// Schema applies the given schema to the database.
func (m Migrator) Schema(schema string) error {
	_, err := m.db.Exec(schema)
	return err
}

// Init ensure the database is initialise for use with the migrator (i.e.
// creates the migrations table).
func (m Migrator) Init() error {
	schema := `
CREATE TABLE IF NOT EXISTS migrations (
  version        INT        PRIMARY KEY
  , checksum     CHAR(64)
  , date_created TIMESTAMP
  , date_applied TIMESTAMP
  , description  TEXT
);
`
	_, err := m.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialise database for migrator: %w", err)
	}
	return nil
}

// Apply applies the given migrations as required. Only migrations that were
// created after the most recently applied migrations was created are applied.
// That is, migrations are only applied if:
//     `migrations.DateCreated > latestMigration.DateCreated`
func (m Migrator) Apply(migrations ...Migration) error {
	err := m.Init()
	if err != nil {
		return err
	}

	latestMigration, err := m.latestMigration()
	if err != nil {
		return err
	}

	m.logger.Debug("Retrieved most recently applied migration", zap.Any("migration", latestMigration))

	migrationsToApply := migrationsAfter(latestMigration, migrations...)
	m.logger.Info("Applying migrations", zap.Int("count", len(migrationsToApply)))

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to get transaction before applying migration: %w", err)
	}
	// Defer the rollback here so we can cleanly exit at any point and have
	// the transaction rollback. If we reach the end of the function and
	// commit the transaction then rollback will do nothing.
	defer tx.Rollback()

	markAppliedStmt := `
INSERT INTO migrations(version, date_created, date_applied, description, checksum)
VALUES($1, $2, now(), $3, $4);
`
	for _, migration := range migrationsToApply {
		m.logger.Info("Applying migration", zap.Int("version", migration.Version), zap.Time("created", migration.Date))
		_, err = tx.Exec(migration.SQL)
		if err != nil {
			return fmt.Errorf("failed to apply migration version %d: %w", migration.Version, err)
		}

		hash := sha256.Sum256([]byte(migration.SQL))
		checksum := fmt.Sprintf("%x", hash)
		_, err := tx.Exec(markAppliedStmt, migration.Version, migration.Date, migration.Description, checksum)
		if err != nil {
			return fmt.Errorf("failed to mark migration version %d as applied: %w", migration.Version, err)
		}
	}
	// Now that we've applied all the required migrations successfully, we
	// can commit the transaction.
	tx.Commit()
	return nil
}

func (m Migrator) latestMigration() (*AppliedMigration, error) {
	query := `
SELECT version, date_created, date_applied, description, checksum FROM migrations
ORDER BY version DESC LIMIT 1;
`
	row := m.db.QueryRow(query)
	var latest AppliedMigration
	err := row.Scan(
		&latest.Version,
		&latest.DateCreated,
		&latest.DateApplied,
		&latest.Description,
		&latest.Checksum,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to find the latest applied migration: %w", err)
	}
	return &latest, nil
}

func migrationsAfter(latestMigration *AppliedMigration, migrations ...Migration) []Migration {
	if latestMigration == nil {
		return migrations
	}

	var filteredMigrations []Migration
	for _, migration := range migrations {
		if migration.Version > latestMigration.Version {
			filteredMigrations = append(filteredMigrations, migration)
		}
	}
	return filteredMigrations
}
