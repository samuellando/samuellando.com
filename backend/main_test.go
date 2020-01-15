
package main_test

import (
  "os"
  "testing"

  "."
  "log"
)

var a main.App

func TestMain(m *testing.M) {
  a = main.App{}
  a.Initialize()

  ensureTableExists()

  code := m.Run()

  clearTable()

  os.Exit(code)
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS testPages
(
  id SERIAL,
  title TEXT NOT NULL,
  text NUMERIC(10,2) NOT NULL DEFAULT 0.00,
  PRIMARY KEY (id)
)
`

func ensureTableExists() {
  if _, err := a.DB.Exec(tableCreationQuery); err != nil {
    log.Fatal(err)
  }
}

func clearTable() {
  a.DB.Exec("DELETE FROM testPages")
  a.DB.Exec("ALTER SEQUENCE pages_id_seq RESTART WITH 1")
}
