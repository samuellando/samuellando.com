// This package provides some utility functions for working with databe/sql.
// Primarly for connecting, and dealing with migrations
package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Options that can be provided to [Connect()]
type Options struct {
	RetrySecs     int    // How many secounds to wait before retrying a failed connection (default: 10)
	MigrationsDir string // The directory to look for migraitons (default: "./migrations")
}

// Functional option setter, must be provided to [Connect()] to modify defaults
type Option func(*Options)

// Connect to a postgres database with the specified information.
//
// Will retry every [Options.RetrySecs] until it succeeds. If the value is configured
// to be less than 0, it will panic on a fialed connection.
//
// After a connection is established, the function will internally call
// [InitalizeMigrations] and [ValidateMigrations], and panic if that fails.
// Migration checking can be disabled by providing an empty string.
func ConnectPostgres(host, port, user, password, dbName string, opts ...Option) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	return connect("postgres", psqlInfo, opts...)
}

// Connect to a sqlite database with the specified information.
//
// Works identically to [ConnectPostgres]
func ConnectSQLite(file string, opts ...Option) *sql.DB {
	return connect("sqlite3", file, opts...)
}

func connect(dbType, info string, opts ...Option) *sql.DB {
	options := Options{
		RetrySecs:     10,
		MigrationsDir: "./migrations",
	}

	for _, opt := range opts {
		opt(&options)
	}

	var db *sql.DB
	var err error
	for {
		if db, err = sql.Open(dbType, info); err == nil {
			if err = db.Ping(); err == nil {
				break
			}
		} else {
			if options.RetrySecs < 0 {
				panic(err)
			}
			log.Println(err)
			log.Printf("Will retry to connect in %d secounds\n", options.RetrySecs)
			time.Sleep(time.Second * time.Duration(options.RetrySecs))
		}
	}
	log.Println("Succesfully connected to database")
	for {
		log.Println("Ensuring base schema")
		err = InitalizeMigrations(db)
        if err == nil {
            break
        }
		if options.RetrySecs < 0 {
			panic(err)
		}
		log.Println(err)
		log.Printf("Will retry to setup schema in %d secounds\n", options.RetrySecs)
		time.Sleep(time.Second * time.Duration(options.RetrySecs))
	}
	if options.MigrationsDir != "" {
		for {
			log.Println("Validating migrations")
			err = ValidateMigrations(db, options.MigrationsDir)
            if err == nil {
                break
            }
			if options.RetrySecs < 0 {
				panic(err)
			}
			log.Println(err)
			log.Printf("Will retry to validate in %d secounds\n", options.RetrySecs)
			time.Sleep(time.Second * time.Duration(options.RetrySecs))
		}
	}
	return db
}
