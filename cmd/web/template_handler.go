package main

import (
	"bufio"
	"fmt"
	"io"
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
	applyFilters(ctxt)
	switch req.Method {
	case "GET":
		if strings.HasPrefix(ctxt.Page, "download") {
			h.downloadDocument(ctxt, w, req)
		} else {
			h.getTemplate(ctxt, w, req)
		}
	case "POST":
		h.createDocument(ctxt, w, req)
	case "PUT":
		h.update(ctxt, w, req)
	case "DELETE":
		h.deleteDocument(ctxt, w, req)
	}
}

func (ctxt *context) GetDocument() *document.Document {
	ds := ctxt.DocumentStore
	doc, err := ds.GetById(ctxt.Reference)
	if err != nil {
		return nil
	}
	return doc
}

func (ctxt *context) GetProject() *project.Project {
	ps := ctxt.ProjectStore
	proj, err := ps.GetById(ctxt.Reference)
	if err != nil {
		return nil
	}
	return proj
}

func (h *templateHandler) getTemplate(ctxt *context, w http.ResponseWriter, req *http.Request) {
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

func (h *templateHandler) downloadDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
	doc := ctxt.GetDocument()
	filename := fmt.Sprintf("%s.md", doc.Title())
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "text/markdown")
	w.WriteHeader(http.StatusOK)
	content, err := doc.Content()
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *templateHandler) createDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
	title := req.PostFormValue("title")
	content := req.PostFormValue("content")
	tags := strings.Split(req.PostFormValue("tags"), ",")
	doc := document.CreateProto(func(df *document.DocumentFeilds) {
		df.Title = title
		df.Content = content
		df.Tags = tags
	})
	err := h.DocumentStore.Add(doc)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	http.Redirect(w, req, fmt.Sprintf("%s/%d", ctxt.Page, doc.Id()), 303)
}

func (h *templateHandler) update(ctxt *context, w http.ResponseWriter, req *http.Request) {
	switch ctxt.Page {
	case "project":
		h.updateProject(ctxt, w, req)
	default:
		h.updateDocument(ctxt, w, req)
	}
}

func (h *templateHandler) updateProject(ctxt *context, w http.ResponseWriter, req *http.Request) {
	proj := ctxt.GetProject()
	desc := req.PostFormValue("description")
	tags := strings.Split(req.PostFormValue("tags"), ",")
	err := proj.Update(func(pf *project.ProjectFields) {
		pf.Description = desc
		pf.Tags = tags
	})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	// And return the updated document
	err = h.templates.ExecuteTemplate(w, "project-li", []any{proj, ctxt})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
}

func (h *templateHandler) updateDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
	doc := ctxt.GetDocument()
	title := req.PostFormValue("title")
	tags := strings.Split(req.PostFormValue("tags"), ",")
	content, err, err_code := getUploadContent(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), err_code)
		return
	}
	if content == "" {
		content = req.PostFormValue("content")
	}
	err = doc.Update(func(df *document.DocumentFeilds) {
		df.Title = title
		df.Content = content
		df.Tags = tags
	})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	// And return the updated document
	err = h.templates.ExecuteTemplate(w, "document", doc)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
}

func (h *templateHandler) deleteDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
	err := h.DocumentStore.Remove(ctxt.GetDocument())
	// Failed to delete
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	http.Redirect(w, req, "/"+ctxt.Page, 303)
	return
}

// Returns the contents of the uploaded file int the "file" field.
// Returns an empty string if there is no file
// Returns a non nil error if there are any problems, along with a status code.
func getUploadContent(req *http.Request) (string, error, int) {
	const max_file_size = int64(2000000)
	f, header, err := req.FormFile("file")
	// If there is no file
	if err != nil {
		return "", nil, 0
	}
	defer f.Close()
	if header.Size > max_file_size {
		return "", fmt.Errorf("%s : %s", http.StatusText(413), "File too large (2MB max)"), 413
	}
	buff := make([]byte, header.Size)
	for {
		r := bufio.NewReader(f)
		_, err = r.Read(buff)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("%s : %s", http.StatusText(500), err), 500
		}
		if err != io.EOF {
			break
		}
	}
	return string(buff), nil, 0
}

func applyFilters(c *context) {
	group := c.Req.FormValue("group")
	tagFiltering := c.Req.FormValue("filter-out-tags")
	tags := make([]string, 0)
	if tagFiltering == "true" {
		if arr, ok := c.Req.Form["filter-tag"]; ok {
			for _, v := range arr {
				tags = append(tags, v)
			}
		}
	}
	// And grab the associated function
	switch c.Page {
	case "projects":
		if tagFiltering == "true" {
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
		sort := "byLastPush"
		if group == "" {
			group = "byLastPush"
			c.Req.Form.Add("group", "byLastPush")
		}
		if sortFunc, ok := c.ProjectSortFunctions[sort]; ok {
			c.ProjectStore = c.ProjectStore.Sort(sortFunc.Func).(*project.Store)
		}
		groupFunc, ok := c.ProjectGroupFunctions[group]
		if ok {
			c.ProjectGroups = c.ProjectStore.Group(groupFunc.Func)
		}
	default:
		if tagFiltering == "true" {
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
		sort := "byCreated"
		if sortFunc, ok := c.DocumentSortFunctions[sort]; ok {
			c.DocumentStore = c.DocumentStore.Sort(sortFunc.Func).(*document.Store)
		}
	}
}
