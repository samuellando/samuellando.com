package db

import (
	"database/sql"
	"fmt"
    "log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Setup the base schema reuqired for managing migrations.
func InitalizeMigrations(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        migration_name TEXT NOT NULL UNIQUE,
        applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err = tx.Exec(query); err != nil {
		return fmt.Errorf("Failed to exec query: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return nil
}

// Go through all the sql scrips in the specified directory and ensure the
// Migrations have been applied
func ValidateMigrations(db *sql.DB, dir string) error {
	missing, err := getMissingMigrations(db, dir)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return fmt.Errorf("Migration %s is missing", missing[0].Name())
	}
	return nil
}

// Go through all the sql scrips in the specified directory and apply the
// Ones that have not yet been applied, in lexical order.
func ApplyMigrations(db *sql.DB, dir string) error {
	missing, err := getMissingMigrations(db, dir)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, file := range missing {
        log.Println("Applying:", file.Name())
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}
		query := string(b)
		if _, err := tx.Exec(query); err != nil {
			return err
		}
		successQuery := `INSERT INTO migrations (migration_name) VALUES ($1)`
		if _, err := tx.Exec(successQuery, file.Name()); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

type byFileName []os.DirEntry

func (a byFileName) Len() int           { return len(a) }
func (a byFileName) Less(i, j int) bool { return strings.Compare(a[i].Name(), a[j].Name()) < 0 }
func (a byFileName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func getMissingMigrations(db *sql.DB, dir string) ([]os.DirEntry, error) {
	// Read migration files
	files, err := os.ReadDir(dir)
	if err != nil {
        return nil, err
	}
    // Verify that there is at least one migration file.
	if len(files) == 0 {
		return nil, fmt.Errorf("No migration files found in %s", dir)
	}
    // Sort the files in order and list the missing ones.
	missing := make([]os.DirEntry, 0, len(files))
	sort.Sort(byFileName(files))
	for _, file := range files {
		migrationName := file.Name()
		var count int
		query := `SELECT COUNT(*) FROM migrations WHERE migration_name = $1`
		err := db.QueryRow(query, migrationName).Scan(&count)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			missing = append(missing, file)
		}
	}
	return missing, nil
}
