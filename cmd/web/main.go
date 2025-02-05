package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"samuellando.com/internal/db"
	"samuellando.com/internal/middleware"
	"samuellando.com/internal/search"
	"samuellando.com/internal/store/document"
	"samuellando.com/internal/store/project"
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
	templates, err := template.New("templates").Funcs(template.FuncMap{
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
	}).ParseGlob(TEMPLATE_DIR + "/*")
	if err != nil {
		panic(err)
	}
	db := db.ConnectPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	defer db.Close()
	documentStore := document.CreateStore(db)
	projectStore := project.CreateStore(db)

	th := templateHandler{
		templates:     *templates,
		DocumentStore: documentStore,
		ProjectStore:  projectStore,
	}

	// Handling static assets
	asset_hander := http.StripPrefix(ASSETS_PREFIX, http.FileServer(http.Dir(ASSETS_DIR)))
	http.Handle(fmt.Sprintf("GET %s/{asset}", ASSETS_PREFIX), asset_hander)
	// Authentication endpoints
	http.HandleFunc("POST /auth", middleware.LoggingFunc(authenticate))
	http.HandleFunc("POST /deauth", middleware.LoggingFunc(deauthenticate))
    // Search endpoint
	http.HandleFunc("GET /search", middleware.LoggingFunc(createSearchHandler(
        search.GenerateIndex("Document", "/writing", &documentStore),
        search.GenerateIndex("Project", "/project", &projectStore),
    )))
	// Template and document CRUD handlers
	http.Handle("/{$}", middleware.Logging(&th))
	http.Handle("/{template}", middleware.Logging(&th))
	http.Handle("/{template}/{document}", middleware.Logging(&th))

	http.ListenAndServe(":8080", nil)
}

func createSearchHandler(indexes ...func() []search.IndexItem) http.HandlerFunc {
    tmpl := `
    <div>
    {{.Type}} <a href="{{.Path}}">{{.Item.Title}}</a>
    </div>
    `
    t, err := template.New("result").Parse(tmpl)
    if err != nil {
        panic(err)
    }
	return func(w http.ResponseWriter, req *http.Request) {
		searchString := req.FormValue("q")
		if searchString == "" {
			return
		}
		for _, index := range indexes {
			elements := index()
			for _, elem := range elements {
				if strings.Contains(elem.Text, searchString) {
                    t.Execute(w, elem)
				}
			}
		}
	}
}
