package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"samuellando.com/internal/auth"
	"samuellando.com/internal/db"
	"samuellando.com/internal/middleware"
	"samuellando.com/internal/search"
	"samuellando.com/internal/store/asset"
	"samuellando.com/internal/store/document"
	"samuellando.com/internal/store/project"
)

const TEMPLATE_DIR = "./templates"
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
	templates := parseTemplates()
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
	}
	ah := assetHandler{
		Store:     &assetStore,
		Templates: *templates,
	}

	// Handling static assets
	static_hander := http.StripPrefix(STATIC_PREFIX, http.FileServer(http.Dir(STATIC_DIR)))
	http.Handle(fmt.Sprintf("GET %s/{asset}", STATIC_PREFIX), static_hander)
	// Authentication endpoints
	http.HandleFunc("POST /auth", middleware.LoggingFunc(auth.Authenticate))
	http.HandleFunc("POST /deauth", middleware.LoggingFunc(auth.Deauthenticate))
	// Search endpoint
	http.HandleFunc("GET /search", middleware.LoggingFunc(search.CreateSearchHandler(
		search.GenerateIndex("Document", "/writing", &documentStore),
		search.GenerateIndex("Project", "/project", &projectStore),
	)))
	// Template and document CRUD handlers
	http.Handle("GET /{$}", middleware.Logging(&th))
	http.Handle("GET /{template}", middleware.Logging(&th))
	http.Handle("GET /{template}/{document}", middleware.Logging(&th))
	// Authenticated actions
	http.Handle("/{$}", middleware.Logging(middleware.Authenticated(&th)))
	http.Handle("/{template}", middleware.Logging(middleware.Authenticated(&th)))
	http.Handle("/{template}/{document}", middleware.Logging(middleware.Authenticated(&th)))
	// Handling user assets
	http.Handle("GET /asset", middleware.Logging(&ah))
	http.Handle("POST /asset", middleware.Logging(middleware.Authenticated(&ah)))
	http.Handle("GET /asset/{asset}", middleware.Logging(&ah))
	http.Handle("DELETE /asset/{asset}", middleware.Logging(middleware.Authenticated(&ah)))

	http.ListenAndServe(":8080", nil)
}

func parseTemplates() *template.Template {
	templ := template.New("templates").Funcs(template.FuncMap{
		"join": strings.Join,
		"byTag": func(needs string) func(*document.Document) bool {
			return func(d *document.Document) bool {
				for _, tag := range d.Tags() {
					if tag == needs {
						return true
					}
				}
				return false
			}
		},
		"arr": func(els ...any) []any {
			return els
		},
		"includes": func(s string, arr []string) bool {
			for _, v := range arr {
				if v == s {
					return true
				}
			}
			return false
		},
	})
	err := filepath.Walk(TEMPLATE_DIR, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return templ
}
