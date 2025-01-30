package main

import (
    "samuellando.com/internal/db"
    "os"
)

var (
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

func main() {
    con := db.ConnectPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, func(opts *db.Options) {
        opts.MigrationsDir = ""
    })
    if err := db.ApplyMigrations(con); err != nil {
        panic(err)
    }
}
