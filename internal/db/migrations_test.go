package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setUp() (*sql.DB, string) {
    db := ConnectSQLite(":memory:", func(opts *Options) {
        opts.RetrySecs = -1
        opts.MigrationsDir = ""
    })
	if db == nil {
		panic("Failed in memory test DB")
	}
	dir, err := os.MkdirTemp("", "db_tests")
	if err != nil {
		panic(err)
	}
	err = InitalizeMigrations(db)
	if err != nil {
		panic(err)
	}
	return db, dir
}

func tearDown(db *sql.DB, dir string) {
	db.Close()
	os.RemoveAll(dir)
}

func TestInitializeMigrations(t *testing.T) {
    db := ConnectSQLite(":memory:", func(opts *Options) {
        opts.RetrySecs = -1
        opts.MigrationsDir = ""
    })
	if db == nil {
		panic("Failed in memory test DB")
	}
    defer db.Close()
    InitalizeMigrations(db)
	migrationName := "0001_example"
	query := `INSERT INTO migrations (migration_name) VALUES ($1)`
    _, err := db.Exec(query, migrationName)
	if err != nil {
		t.Error("The migrations table should exist")
	}
}

func TestValidateMigrationsNoFiles(t *testing.T) {
	db, dir := setUp()
	defer tearDown(db, dir)
	if err := ValidateMigrations(db, dir); err == nil {
		t.Error("Should return an error when there are no migration files")
	}
}

func TestValidateMigrationsMissing(t *testing.T) {
	db, dir := setUp()
	defer tearDown(db, dir)
	file, err := os.Create(filepath.Join(dir, "001_example.sql"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := ValidateMigrations(db, dir); err == nil {
		t.Error("Should return an error when a migration has not been applied")
	}
}

func TestValidateMigrationsOkay(t *testing.T) {
	db, dir := setUp()
	defer tearDown(db, dir)
	file_name := "0001_example.sql"
	query := `INSERT INTO migrations (migration_name) VALUES ($1);`
	if _, err := db.Exec(query, file_name); err != nil {
		panic(err)
	}
	file, err := os.Create(filepath.Join(dir, file_name))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(query))

	if err := ValidateMigrations(db, dir); err != nil {
		t.Error("Should not fail")
	}
}

func TestApplyMigrationsNoFiles(t *testing.T) {
	db, dir := setUp()
	defer tearDown(db, dir)
	if err := ApplyMigrations(db, func(o *Options) {
        o.MigrationsDir = dir
    }); err == nil {
		t.Error("Should return an error when there are no migration files")
	}
}

func TestAppMigrationsUpdatesTable(t *testing.T) {
	db, dir := setUp()
	defer tearDown(db, dir)

	file_name := "0001_example.sql"
	query := `CREATE TABLE example (id int DEFAULT 1);`

	err := os.WriteFile(filepath.Join(dir, file_name), []byte(query), 0644)
	if err != nil {
		panic(err)
	}

	if err := ApplyMigrations(db, func(o *Options) {
        o.MigrationsDir = dir
    }); err != nil {
		t.Error("Should not fail")
	}
    query = "SELECT count(*) FROM migrations;"
    var count int
    row := db.QueryRow(query)
    row.Scan(&count)
    if count == 0 {
        t.Error("Should update the table")
    }
}
func TestAppMigrationsOrder(t *testing.T) {
	db, dir := setUp()
	defer tearDown(db, dir)

	file_name1 := "0001_example.sql"
	query1 := `CREATE TABLE example (id int DEFAULT 1);`
	file_name2 := "0002_example.sql"
	query2 := `ALTER TABLE example ADD id2 int DEFAULT 2;`

	err := os.WriteFile(filepath.Join(dir, file_name1), []byte(query1), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(dir, file_name2), []byte(query2), 0644)
	if err != nil {
		panic(err)
	}

	if err := ApplyMigrations(db, func(o *Options) {
        o.MigrationsDir = dir
    }); err != nil {
		t.Error("Should not fail")
	}
	query := `INSERT INTO example (id) VALUES (123)`
	_, err = db.Exec(query)
	if err != nil {
		t.Errorf("Should be able to isnert the first migration: %s\n", err)
	}
	query = `INSERT INTO example (id2) VALUES (345)`
	_, err = db.Exec(query)
	if err != nil {
		t.Errorf("Should be able to isnert the second migration: %s\n", err)
	}
}
