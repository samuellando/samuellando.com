package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"samuellando.com/internal/auth"
	"samuellando.com/internal/db"
	"samuellando.com/internal/middleware"
	"samuellando.com/internal/search"
	"samuellando.com/internal/store/asset"
	"samuellando.com/internal/store/document"
	"samuellando.com/internal/store/project"
	"samuellando.com/internal/template"
)

const TEMPLATE_DIR = "templates"
const STATIC_DIR = "./static"
const STATIC_PREFIX = "/static"

var (
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

func main() {
	templates := template.New("templates").Funcs(template.FuncMap{
		"join": strings.Join,
		"byTag": func(needs string) func(*document.Document) bool {
			return func(d *document.Document) bool {
				return slices.Contains(d.Tags(), needs)
			}
		},
		"arr": func(els ...any) []any {
			return els
		},
		"includes": func(s string, arr []string) bool {
			return slices.Contains(arr, s)
		},
	}).ParseTemplates(TEMPLATE_DIR)
	db := db.ConnectPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
		func(o *db.Options) {
			o.RetrySecs = -1
		})
	defer db.Close()
	documentStore := document.CreateStore(db)
	projectStore := project.CreateStore(db)
	assetStore := asset.CreateStore(db)

	th := templateHandler{
		templates:     *templates,
		DocumentStore: documentStore,
		ProjectStore:  projectStore,
		AssetStore:    assetStore,
	}
	ah := assetHandler{
		Store:     &assetStore,
		Templates: *templates,
	}
	dh := documentHandler{
		templates:     *templates,
		documentStore: documentStore,
	}
	ph := projectHandler{
		templates:    *templates,
		projectStore: projectStore,
	}

	// Handling static assets
	static_hander := http.StripPrefix(STATIC_PREFIX, http.FileServer(http.Dir(STATIC_DIR)))
	http.Handle(fmt.Sprintf("GET %s/{asset}", STATIC_PREFIX), static_hander)
	// Authentication endpoints
	http.HandleFunc("POST /auth", middleware.LoggingFunc(auth.Authenticate))
	http.HandleFunc("POST /deauth", middleware.LoggingFunc(auth.Deauthenticate))
	// Search endpoint
	searchResultTemplate := templates.Lookup("search-result")
	if searchResultTemplate == nil {
		panic("Must define search result template")
	}
	http.HandleFunc("GET /search", middleware.LoggingFunc(
		search.CreateSearchHandler(
			*searchResultTemplate,
			search.GenerateIndex("Project", "/project", &projectStore),
		)))
	// Template
	http.Handle("GET /", middleware.Logging(&th))
	// Document actions
	http.Handle("GET /document/{document}", middleware.Logging(&dh))
	http.Handle("POST /document", middleware.Logging(middleware.Authenticated(&dh)))
	http.Handle("PUT /document/{document}", middleware.Logging(middleware.Authenticated(&dh)))
	http.Handle("DELETE /document/{document}", middleware.Logging(middleware.Authenticated(&dh)))
	// Project actions
	http.Handle("GET /project/{project}", middleware.Logging(&ph))
	http.Handle("PUT /project/{project}", middleware.Logging(middleware.Authenticated(&ph)))
	// Handling user assets
	http.Handle("GET /asset", middleware.Logging(&ah))
	http.Handle("POST /asset", middleware.Logging(middleware.Authenticated(&ah)))
	http.Handle("GET /asset/{asset}", middleware.Logging(&ah))
	http.Handle("DELETE /asset/{asset}", middleware.Logging(middleware.Authenticated(&ah)))

	http.ListenAndServe(":8080", nil)
}
