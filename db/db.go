package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type ConnectionDetails struct {
	User     string
	Password string
	DBName   string
	Host     string
	SSLMode  string
}

func NewConnectionDetails(user, password, dbName, host, sslmode string) ConnectionDetails {
	return ConnectionDetails{
		User:     strings.TrimSpace(user),
		Password: strings.TrimSpace(password),
		DBName:   strings.TrimSpace(dbName),
		Host:     strings.TrimSpace(host),
		SSLMode:  strings.TrimSpace(sslmode),
	}
}

func (c ConnectionDetails) String() string {
	mapping := make(map[string]string)
	if c.User != "" {
		mapping["user"] = c.User
	}
	if c.Password != "" {
		mapping["password"] = c.Password
	}
	if c.DBName != "" {
		mapping["dbname"] = c.DBName
	}
	if c.Host != "" {
		mapping["host"] = c.Host
	}
	if c.SSLMode != "" {
		mapping["sslmode"] = c.SSLMode
	}

	var connString strings.Builder
	prefix := ""
	for key, value := range mapping {
		fmt.Fprintf(&connString, "%s%s=%s", prefix, key, value)
		prefix = " "
	}
	return connString.String()
}

// DBConn returns a database connection (or error if it can't connect) based on
// the given user, password and host. It will attempt to connect up to 20 times
// with an exponential back off.
func NewConn(logger *zap.Logger, user, password, dbName, host, sslmode string) (*sqlx.DB, error) {
	connDetails := NewConnectionDetails(user, password, dbName, host, sslmode)
	connString := connDetails.String()

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, xerrors.Errorf("failed to open database: %w", err)
	}

	logger.Info("Attempting to connect to database",
		zap.String("host", host), zap.String("dbName", dbName), zap.String("user", user),
	)
	maxRetries := 20
	for retry := 1; retry <= maxRetries; retry++ {
		if err = db.Ping(); err == nil {
			// If we can connect to the db okay then there is no point retrying
			// anymore so just exit here.
			return db, nil
		}
		logger.Info("Failed to connect to database",
			zap.Int("attempt", retry), zap.String("dbName", dbName), zap.String("host", host), zap.String("user", user), zap.Error(err),
		)
		// Exponentially back off so we don't spam the db too much.
		time.Sleep(time.Duration(retry) * time.Second)
	}
	return nil, xerrors.Errorf("failed to connect to database %s after %d retries: %w", dbName, maxRetries, err)
}
