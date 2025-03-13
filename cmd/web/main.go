package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	htmlTemplate "html/template"
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

var TEMPLATE_FUNCTIONS = template.FuncMap{
	"join": strings.Join,
	"joinTags": func(tags []tag.ProtoTag, sep string) string {
		w := new(strings.Builder)
		for i, t := range tags {
			w.WriteString(t.Value)
			if i < len(tags)-1 {
				w.WriteString(sep)
			}
		}
		return w.String()
	},
	"byTag": func(needs string) func(document.Document) bool {
		return func(d document.Document) bool {
			for _, t := range d.Tags() {
				if t.Value == needs {
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
}

func main() {
	templates := template.New("templates").Funcs(TEMPLATE_FUNCTIONS).ParseTemplates(TEMPLATE_DIR)
	db := db.ConnectPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
		func(o *db.Options) {
			// Crash as soon as possible in case of a db connection/migrations issue.
			o.RetrySecs = -1
		})
	defer db.Close()

	documentStore := document.CreateStore(db)
	projectStore := project.CreateStore(db)
	assetStore := asset.CreateStore(db)
	tagStore := tag.CreateStore(db)

	th := template.Handler{
		Templates: *templates,
		// At this point, we throw away type safety for convinience on the frontend.
		ContextValues: template.ContextValues{
			"DocumentStore": func(ctx template.Context) any { return documentStore },
			"ProjectStore":  func(ctx template.Context) any { return projectStore },
			"ProjectGroups": func(ctx template.Context) any {
				filterTags := ctx.Get("FilterTags").([]string)
				filtered, err := projectStore.Filter(func(p project.Project) bool {
					if len(filterTags) == 0 {
						return true
					}
					for _, t := range filterTags {
						for _, pt := range p.Tags() {
							if pt.Value == t {
								return true
							}
						}
					}
					return false
				})
				if err != nil {
					filtered = projectStore
				}
				sorted, err := filtered.Sort(func(p1, p2 project.Project) bool {
					return p1.Pushed().After(p2.Pushed())
				})
				if err != nil {
					sorted = filtered
				}
				groups, err := sorted.Group(func(p project.Project) string {
					return p.Pushed().Format("2006")
				})
				if err != nil {
					return nil
				}
				return groups
			},
			"AssetStore": func(ctx template.Context) any { return assetStore },
			"TagStore":   func(ctx template.Context) any { return tagStore },
			"Admin": func(ctx template.Context) any {
				return auth.IsAuthenticated(ctx.Get("Req").(*http.Request))
			},
			"Reference": func(ctx template.Context) any {
				path := ctx.Get("Page").(string)
				parts := strings.Split(path, "/")
				ref := parts[len(parts)-1]
				id, err := strconv.Atoi(ref)
				if err != nil {
					return -1
				}
				return id
			},
			"Document": func(ctx template.Context) any {
				id := ctx.Get("Reference").(int)
				doc, err := documentStore.GetById(int64(id))
				if err != nil {
					return nil
				}
				return doc
			},
			"Project": func(ctx template.Context) any {
				id := ctx.Get("Reference").(int)
				proj, err := projectStore.GetById(int64(id))
				if err != nil {
					return nil
				}
				return proj
			},
			"FilterTags": func(ctx template.Context) any {
				req := ctx.Get("Req").(*http.Request)
				err := req.ParseForm()
				if err != nil {
					return []string{}
				}
				if f, ok := req.Form["filter-tag"]; ok {
					return f
				} else {
					return []string{}
				}
			},
		},
	}

	ah := asset.Handler{
		Store: assetStore,
	}
	tagh := tag.Handler{
		Store: tagStore,
	}

	docTemplate := templates.Lookup("document")
	if docTemplate == nil {
		panic("Must define document template")
	}
	dh := document.Handler{
		Template:      *docTemplate,
		DocumentStore: documentStore,
		TagStore:      tagStore,
	}

	projTemplate := templates.Lookup("project")
	if projTemplate == nil {
		panic("Must define project template")
	}
	ph := project.Handler{
		Template:     *projTemplate,
		ProjectStore: projectStore,
		TagStore:     tagStore,
	}

	searchResultTemplate := templates.Lookup("search-result")
	if searchResultTemplate == nil {
		panic("Must define search result template")
	}
	sh := createSearchHandler(searchResultTemplate, projectStore)

	// Handling static assets
	static_hander := http.StripPrefix(STATIC_PREFIX, http.FileServer(http.Dir(STATIC_DIR)))
	http.Handle(fmt.Sprintf("GET %s/{asset}", STATIC_PREFIX), static_hander)
	// Authentication endpoints
	http.HandleFunc("POST /auth", middleware.LoggingFunc(auth.Authenticate))
	http.HandleFunc("POST /deauth", middleware.LoggingFunc(auth.Deauthenticate))
	// Search endpoint
	http.HandleFunc("GET /search", middleware.LoggingFunc(sh))
	// In general, everything should get served by the template handler.
	http.Handle("GET /", middleware.Logging(&th))
	// Handling user assets
	http.Handle("GET /asset/{asset}", middleware.Logging(&ah))
	http.Handle("POST /asset", middleware.Logging(middleware.Authenticated(&ah)))
	http.Handle("DELETE /asset/{asset}", middleware.Logging(middleware.Authenticated(&ah)))
	// Document actions
	http.Handle("POST /document", middleware.Logging(middleware.Authenticated(&dh)))
	http.Handle("PUT /document/{document}", middleware.Logging(middleware.Authenticated(&dh)))
	http.Handle("DELETE /document/{document}", middleware.Logging(middleware.Authenticated(&dh)))
	// Project actions
	http.Handle("PUT /project/{project}", middleware.Logging(middleware.Authenticated(&ph)))
	// Handling tag edits
	http.Handle("PATCH /tag/{tag}", middleware.Logging(middleware.Authenticated(&tagh)))
	http.Handle("DELETE /tag/{tag}", middleware.Logging(middleware.Authenticated(&tagh)))

	http.ListenAndServe(":8080", nil)
}

func createSearchHandler(template *htmlTemplate.Template, projectStore project.Store) http.HandlerFunc {
	searchStore, err := projectStore.Filter(func(p project.Project) bool {
		return !p.Hidden()
	})
	if err != nil {
		panic(err)
	}
	ps, ok := searchStore.(project.Store)
	if !ok {
		panic("Not a project store, this should not happen")
	}
	return search.CreateSearchHandler(
		*template,
		search.GenerateIndex("Project", "/projects", &ps),
	)
}
