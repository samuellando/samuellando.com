
package main_test

import (
  "os"
  "testing"

  "."
  "log"
  "net/http"
  "net/http/httptest"
  "encoding/json"
  "bytes"
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
CREATE TABLE IF NOT EXISTS pages
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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/page", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentPage(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/page/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Page not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Page not found'. Got '%s'", m["error"])
	}
}

func TestCreatePage(t *testing.T) {
  clearTable()

  type Page struct {
		title string
		text  string
	}
	page := Page {
		title:  "test page",
		text:   "this is a test page",
	}

  payload, _ := json.Marshal(page)

	req, _ := http.NewRequest("POST", "/page", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "test page" {
		t.Errorf("Expected page title to be 'test page'. Got '%v'", m["title"])
	}

	if m["text"] != 11.22 {
		t.Errorf("Expected page text to be 'this is a test page'. Got '%v'", m["text"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected page ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetPage(t *testing.T) {
	clearTable()
	addPages(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addPages(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO pages(title, text) VALUES($1, $2)", "page pagetext")
	}
}

func TestUpdatePage(t *testing.T) {
	clearTable()
	addPages(1)

	req, _ := http.NewRequest("GET", "/page/1", nil)
	response := executeRequest(req)
	var originalPage map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &originalPage)

  type Page struct {
		title string
		text  string
	}
	page := Page {
		title:  "new title",
		text:   "new text",
	}

  payload, _ := json.Marshal(page)

	req, _ = http.NewRequest("PUT", "/page/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalPage["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalPage["id"], m["id"])
	}

	if m["title"] == originalPage["title"] {
		t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalPage["title"], m["title"], m["title"])
	}

	if m["text"] == originalPage["text"] {
		t.Errorf("Expected the text to change from '%v' to '%v'. Got '%v'", originalPage["text"], m["text"], m["text"])
	}
}

func TestDeletePage(t *testing.T) {
	clearTable()
	addPages(1)

	req, _ := http.NewRequest("GET", "/page/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/page/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/page/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}