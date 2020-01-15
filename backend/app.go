package main

import (
  "database/sql"
  "github.com/gorilla/mux"
  _ "github.com/lib/pq"

  "fmt"
  "log"
)

const (
  host     = "localhost"
  user     = "pageAdmin"
  password = "password"
  dbname   = "personalWebsite"
  
)

type App struct {
  Router *mux.Router
  DB     *sql.DB
}

func (a *App) Initialize() {
  sqlInfo := fmt.Sprintf("host=%s user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, user, password, dbname)

  var err error
  a.DB, err = sql.Open("postgres", sqlInfo)
  if err != nil {
    log.Fatal(err)
  }

  a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) { }
