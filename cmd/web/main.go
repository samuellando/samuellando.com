package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"samuellando.com/internal/db"
	"samuellando.com/internal/middleware"
	"samuellando.com/internal/stores"
)

const TEMPLATE_DIR = "./templates"
const ASSETS_DIR = "./assets"
const ASSETS_PREFIX = "/assets"

var (
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

func main() {
	templates, err := template.New("templates").Funcs(template.FuncMap{"join": strings.Join}).ParseGlob(TEMPLATE_DIR + "/*")
	if err != nil {
		panic(err)
	}
	db := db.ConnectPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	defer db.Close()
	markdownStore := stores.InitializeMarkdownStore(db)
	projectStore := stores.InitializeProjectStore()

	th := templateHandler{templates: *templates, MarkdownStore: markdownStore, ProjectStore: projectStore}

	// Handling static assets
	asset_hander := http.StripPrefix(ASSETS_PREFIX, http.FileServer(http.Dir(ASSETS_DIR)))
	http.Handle(fmt.Sprintf("GET %s/{asset}", ASSETS_PREFIX), asset_hander)
	// Authentication endpoints
	http.HandleFunc("POST /auth", middleware.LoggingFunc(authenticate))
	http.HandleFunc("POST /deauth", middleware.LoggingFunc(deauthenticate))
	// Template and document CRUD handlers
	http.Handle("/{$}", middleware.Logging(&th))
	http.Handle("/{template}", middleware.Logging(&th))
	http.Handle("/{template}/{document}", middleware.Logging(&th))

	http.ListenAndServe(":8080", nil)
}
