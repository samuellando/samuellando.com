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
	"samuellando.com/internal/store/tag"
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
		"joinTags": func(tags []tag.Tag, sep string) string {
			w := new(strings.Builder)
			for i, t := range tags {
				w.WriteString(t.Value())
				if i < len(tags)-1 {
					w.WriteString(sep)
				}
			}
			return w.String()
		},
		"byTag": func(needs string) func(*document.Document) bool {
			return func(d *document.Document) bool {
				for _, t := range d.Tags() {
					if t.Value() == needs {
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
	tagStore := tag.CreateStore(db)

	th := templateHandler{
		templates:     *templates,
		DocumentStore: documentStore,
		ProjectStore:  projectStore,
		AssetStore:    assetStore,
		TagStore:      tagStore,
	}
	ah := assetHandler{
		Store:     &assetStore,
		Templates: *templates,
	}
	tagh := tagHandler{
		Store: &tagStore,
	}
	dh := documentHandler{
		templates:     *templates,
		documentStore: documentStore,
		tagStore:      tagStore,
	}
	ph := projectHandler{
		templates:    *templates,
		projectStore: projectStore,
		tagStore:     tagStore,
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
	// Handling tag edits
	http.Handle("PATCH /tag/{tag}", middleware.Logging(middleware.Authenticated(&tagh)))
	http.Handle("DELETE /tag/{tag}", middleware.Logging(middleware.Authenticated(&tagh)))

	http.ListenAndServe(":8080", nil)
}
