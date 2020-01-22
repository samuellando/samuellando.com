package main

import (
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"golang.org/x/crypto/bcrypt"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	ensureTablesExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS pages
(
  id SERIAL,
  title TEXT NOT NULL,
  text TEXT NOT NULL,
  private boolean DEFAULT true,
  PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS users
(
  id SERIAL,
  username TEXT UNIQUE NOT NULL ,
  password TEXT NOT NULL,
  apikey TEXT UNIQUE NOT NULL,
  PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS relations
(
  id SERIAL,
  userid int NOT NULL,
  pageid int NOT NULL
  level int NOT NULL
  PRIMARY KEY (id)
);
`

/*
 *				HELPERS
 */

func ensureTablesExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTables() {
	a.DB.Exec("DELETE FROM pages")
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("DELETE FROM relations")
	a.DB.Exec("ALTER SEQUENCE pages_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE relations_id_seq RESTART WITH 1")
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

func addPages(count int, private bool) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.QueryRow("INSERT INTO pages(title, text, private) VALUES($1, $2, $3)", "page", "pageText", private)
	}
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		a.DB.QueryRow("INSERT INTO users(username, password, apikey) VALUES($1, $2, $3)", "username", password, i)
	}
}

func addRelation(uid, pid, level int) {
	a.DB.QueryRow("INSERT INTO relations(userid, pageid, level) VALUES($1, $2, $3)", uid, pid, level)
}

/*
 *				USER TESTS
 */

func TestPostUser(t *testing.T) {
	clearTables()

	type User struct {
		Username string
		Password string
	}

	username := "John"
	password := "password123"

	user := User{
		Username: username,
		Password: password,
	}

	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["username"] != username {
		t.Errorf("Expected username to be '%v'. Got '%v'", username, m["username"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}

	if len(m["apikey"]) > 10 {
		t.Errorf("Expected an APIkey, got '%v'", m["apikey"])
	}
}

func TestFailedPostUser(t *testing.T) {
	clearTables()
	addUsers(1)

	type User struct {
		Username string
		Password string
	}

	username := "username"
	password := "password123"

	user := User{
		Username: username,
		Password: password,
	}

	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusConflict, response.Code)
}

func TestGetEmptyUsers(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) >= 1 {
		t.Errorf("Expected an empty set, got '%v'", len(m))
	}
}

func TestGetUsers(t *testing.T) {
	clearTables()
	addUsers(2)

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) != 2 {
		t.Errorf("Expected an two users got '%v'", len(m))
	}

	if m["username"] != "username" {
		t.Errorf("Expected username 'username' got '%v'", m["username"])
	}
}

func TestGetAPIkey(t *testing.T) {
	clearTables()
	addUsers(1)

	type User struct {
		Username string
		Password string
	}

	username := "username"
	password := "password"

	user := User{
		Username: username,
		Password: password,
	}

	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("GET", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["username"] != username {
		t.Errorf("Expected username to be '%v'. Got '%v'", username, m["username"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}

	if m["apikey"] != "1" {
		t.Errorf("Expected an APIkey, got '%v'", m["apikey"])
	}
}

func TestFailedGetAPIkey(t *testing.T) {
	clearTables()
	addUsers(1)

	type User struct {
		Username string
		Password string
	}

	username := "username"
	password := "wrongpassword"

	user := User{
		Username: username,
		Password: password,
	}

	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("GET", "/user", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

/*
 *				PAGE TESTS
 */

func TestGuestGetPublicPage(t *testing.T) {
	clearTables()
	addPages(1, false)

	req, _ := http.NewRequest("GET", "/page/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "title" {
		t.Errorf("Expected page title to be 'title'. Got '%v'", m["title"])
	}

	if m["text"] != "text" {
		t.Errorf("Expected page text to be 'text'. Got '%v'", m["text"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected page id to be 1, got '%v'", m["id"])
	}
}

func TestFailedGuestGetPrivatePage(t *testing.T) {
	clearTables()
	addPages(1, true)

	req, _ := http.NewRequest("GET", "/page/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestUserGetPublicPage(t *testing.T) {
	clearTables()
	addPages(1, false)
	addUsers(1)

	req, _ := http.NewRequest("GET", "/page/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "title" {
		t.Errorf("Expected page title to be 'title'. Got '%v'", m["title"])
	}

	if m["text"] != "text" {
		t.Errorf("Expected page text to be 'text'. Got '%v'", m["text"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected page id to be 1, got '%v'", m["id"])
	}
}

func TestUserGetPrivatePage(t *testing.T) {}

func TestViewerGetPublicPage(t *testing.T) {}

func TestViewerGetPrivatePage(t *testing.T) {}

func TestEditorGetPublicPage(t *testing.T) {}

func TestEditorGetPrivatePage(t *testing.T) {}

func TestOwnerGetPublicPage(t *testing.T) {}

func TestOwnerGetPrivatePage(t *testing.T) {}

func TestGuestPostPublicPage(t *testing.T) {}

func TestGuestPostPrivatePage(t *testing.T) {}

func TestUserPostPublicPage(t *testing.T) {}

func TestUserPostPrivatePage(t *testing.T) {}

func TestFailedGuestPutPublicPage(t *test.T) {}

func TestFailedGuestPutPrivatePage(t *test.T) {}

func TestFailedUserPutPublicPage(t *test.T) {}

func TestFailedUserPutPrivatePage(t *test.T) {}

func TestFailedViewerPutPublicPage(t *test.T) {}

func TestFailedViewerPutPrivatePage(t *test.T) {}

func TestOwnerPutPublicPage(t *test.T) {}

func TestOwnerPutPrivatePage(t *test.T) {}

func TestEditorPutPublicPage(t *test.T) {}

func TestEditorPutPrivatePage(t *test.T) {}

func TestFailedUserDeletePublicPage(t *testing.T) {}

func TestFailedUserDeletePrivatePage(t *testing.T) {}

func TestFailedViewerDeletePublicPage(t *testing.T) {}

func TestFailedViewerDeletePrivatePage(t *testing.T) {}

func TestFailedEditorDeletePublicPage(t *testing.T) {}

func TestFailedEditorDeletePrivatePage(t *testing.T) {}

func TestOwnerGetDeletePublicPage(t *testing.T) {}

func TestOwnerGetDeletePrivatePage(t *testing.T) {}

func TestGuestGetPages(t *testing.T) {}

func TestUserGetPages(t *testing.T) {}
