package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"samuellando.com/internal/auth"
	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"
	"samuellando.com/internal/store/asset"
	"samuellando.com/internal/store/document"
	"samuellando.com/internal/store/project"
	"samuellando.com/internal/template"
)

type context struct {
	ProjectStore          *project.Store
	ProjectGroups         *datatypes.OrderedMap[string, store.Store[*project.Project]]
	DocumentStore         *document.Store
	AssetStore            *asset.Store
	Page                  string
	Reference             int
	Admin                 bool
	Req                   *http.Request
	ProjectSortFunctions  map[string]sortFunctionReference[*project.Project]
	DocumentSortFunctions map[string]sortFunctionReference[*document.Document]
	ProjectGroupFunctions map[string]groupFunctionReference[*project.Project]
	FilterTags            []string
}

func getPathContext(req *http.Request) (string, string, bool) {
	path := req.URL.Path
	parts := strings.Split(path, "/")
	ref := parts[len(parts)-1]
	return path, ref, auth.IsAuthenticated(req)
}

type templateHandler struct {
	templates     template.Template
	DocumentStore document.Store
	ProjectStore  project.Store
	AssetStore    asset.Store
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	page, ref, admin := getPathContext(req)
	id, err := strconv.Atoi(ref)
	if err != nil {
		id = -1
	}
	ctxt := &context{
		DocumentStore:         &h.DocumentStore,
		ProjectStore:          &h.ProjectStore,
		AssetStore:            &h.AssetStore,
		Reference:             id,
		Page:                  page,
		Admin:                 admin,
		Req:                   req,
		ProjectSortFunctions:  PROJECT_SORT_FUNCTIONS,
		DocumentSortFunctions: DOCUMENT_SORT_FUNCTIONS,
		ProjectGroupFunctions: PROJECT_GROUP_FUNCTIONS,
	}
	ctxt.arrangeStores()
	h.renderTemplate(ctxt, w, req)
}

// Allows to get the requested document from within a template
func (ctxt *context) GetDocument() *document.Document {
	ds := ctxt.DocumentStore
	doc, err := ds.GetById(ctxt.Reference)
	if err != nil {
		return nil

	}
	return doc
}

// Allows to get the requested project from within a template
func (ctxt *context) GetProject() *project.Project {
	ps := ctxt.ProjectStore
	proj, err := ps.GetById(ctxt.Reference)
	if err != nil {
		return nil
	}
	return proj
}

func (h *templateHandler) renderTemplate(ctxt *context, w http.ResponseWriter, req *http.Request) {
	template := path.Join("pages", ctxt.Page)
	// Check that the template exists
	if h.templates.Lookup(template) == nil {
		// Check for a slug
		template = filepath.Dir(template) + "/[slug]"
		if h.templates.Lookup(template) == nil {
			http.NotFound(w, req)
			return
		}
	}
	err := h.templates.ExecuteTemplate(w, template, ctxt)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (c *context) arrangeStores() {
    c.Req.ParseForm()
	tags := make([]string, 0)
    log.Println(c.Req.Form)
	if arr, ok := c.Req.Form["filter-tag"]; ok {
		for _, v := range arr {
			tags = append(tags, v)
		}
	}
	c.FilterTags = tags
	c.arrangeProjects("byLastPush", "byLastPush", tags)
	c.arrangeDocuments("byCreated", tags)
}

func (c *context) arrangeProjects(sortRef string, groupRef string, tags []string) {
	// And grab the associated function
	if len(tags) > 0 {
		c.ProjectStore = c.ProjectStore.Filter(func(p *project.Project) bool {
			for _, pt := range p.Tags() {
				for _, t := range tags {
					if pt == t {
						return true
					}
				}
			}
			return false
		}).(*project.Store)
		c.FilterTags = tags
	}
	if groupRef == "" {
		groupRef = "byLastPush"
		c.Req.Form.Add("group", "byLastPush")
	}
	if sortFunc, ok := c.ProjectSortFunctions[sortRef]; ok {
		c.ProjectStore = c.ProjectStore.Sort(sortFunc.Func).(*project.Store)
	}
	groupFunc, ok := c.ProjectGroupFunctions[groupRef]
	if ok {
		c.ProjectGroups = c.ProjectStore.Group(groupFunc.Func)
	}
}

func (c *context) arrangeDocuments(sortRef string, tags []string) {
	if len(tags) > 0 {
		c.DocumentStore = c.DocumentStore.Filter(func(d *document.Document) bool {
			for _, dt := range d.Tags() {
				for _, t := range tags {
					if dt == t {
						return true
					}
				}
			}
			return false
		}).(*document.Store)
		c.FilterTags = tags
	}
	if sortFunc, ok := c.DocumentSortFunctions[sortRef]; ok {
		c.DocumentStore = c.DocumentStore.Sort(sortFunc.Func).(*document.Store)
	}
}
