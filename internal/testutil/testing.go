// This package contains common utility functions for testing
package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"samuellando.com/internal/db"
)

// Pull db Credentials from environment variables. Compatable with db.ConnectPostgress.
//
// db.ConnectPostgres(GetDbCredentials()... opts) *sql.DB {
func GetDbCredentials() (string, string, string, string, string, func(*db.Options)) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	options := func(opts *db.Options) {
		opts.RetrySecs = -1
		opts.MigrationsDir = ""
		opts.Logger = createDiscardLogger()
	}
	return DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, options
}

// Empty the datasbase schema
func ResetDb(con *sql.DB, schema string) error {
	err := clearSchema(con, schema)
	if err != nil {
		return err
	}
	_, err = con.Exec(fmt.Sprintf("SET search_path TO %s;", schema))
	if err != nil {
		return err
	}
	err = db.InitalizeMigrations(con)
	migrations, err := getMigrationsPath()
	if err != nil {
		return err
	}
	err = db.ApplyMigrations(con, func(o *db.Options) {
		o.MigrationsDir = migrations
	})
	return err
}

func clearSchema(con *sql.DB, schema string) error {
	tx, err := con.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE; CREATE SCHEMA %s;", schema, schema))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Resolves the ./migrations dir from the project root
func getMigrationsPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Backup until we find the migrations dir
	var dir string
	for {
		dir = filepath.Join(wd, "migrations")
		if _, err = os.Stat(dir); err == nil {
			break
		}
		parent := filepath.Dir(wd)
		if parent == wd { // Reached root "/"
			return "", fmt.Errorf("Hit root dir")
		}
		wd = parent
	}
	return dir, nil
}

type discardWriter struct{}

// Write implements io.Writer but does nothing.
func (d *discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil // Pretend we wrote everything successfully
}

func createDiscardLogger() *log.Logger {
	return log.New(&discardWriter{}, "", 0)
}
