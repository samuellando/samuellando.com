package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"samuellando.com/internal/auth"
	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"
	"samuellando.com/internal/store/document"
	"samuellando.com/internal/store/project"
)

type context struct {
	ProjectStore          *project.Store
	ProjectGroups         *datatypes.OrderedMap[string, store.Store[*project.Project]]
	DocumentStore         *document.Store
	Page                  string
	Document              *document.Document
	Project               *project.Project
	Admin                 bool
	Req                   *http.Request
	ProjectSortFunctions  map[string]sortFunctionReference[*project.Project]
	DocumentSortFunctions map[string]sortFunctionReference[*document.Document]
	ProjectGroupFunctions map[string]groupFunctionReference[*project.Project]
}

func getPathContext(req *http.Request) (string, string, bool) {
	path := req.URL.Path
	var page string
	var ref string
	if path == "/" {
		page = "index"
	} else {
		page = req.PathValue("template")
		ref = req.PathValue("document")
	}
	return page, ref, isAuthenticated(req)
}

func isAuthenticated(req *http.Request) bool {
	if cookie, err := req.Cookie("session"); err == nil && auth.ValidJWT(cookie.Value) {
		return true
	} else {
		return false
	}
}

type templateHandler struct {
	templates     template.Template
	DocumentStore document.Store
	ProjectStore  project.Store
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	page, ref, admin := getPathContext(req)
	var doc *document.Document
	var proj *project.Project
	switch page {
	case "project":
		if ref != "" {
			projectId, err := strconv.Atoi(ref)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(400), "project reference must be numeric"), 400)
				return
			}
			proj, err = h.ProjectStore.GetById(projectId)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(404), "project not found"), 404)
				return
			}
		}
	default:
		if ref != "" {
			documentId, err := strconv.Atoi(ref)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(400), "document reference must be numeric"), 400)
				return
			}
			doc, err = h.DocumentStore.GetById(documentId)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(404), "document not found"), 404)
				return
			}
		}
	}
	ctxt := &context{
		DocumentStore:         &h.DocumentStore,
		ProjectStore:          &h.ProjectStore,
		Page:                  page,
		Document:              doc,
		Project:               proj,
		Admin:                 admin,
		Req:                   req,
		ProjectSortFunctions:  PROJECT_SORT_FUNCTIONS,
		DocumentSortFunctions: DOCUMENT_SORT_FUNCTIONS,
		ProjectGroupFunctions: PROJECT_GROUP_FUNCTIONS,
	}
	applyFilters(ctxt)
	switch req.Method {
	case "GET":
		h.getTemplate(ctxt, w, req)
	case "POST":
		h.createDocument(w, req)
	case "PUT":
		h.update(ctxt, w, req)
	case "DELETE":
		h.deleteDocument(ctxt, w, req)
	}
}

func (h *templateHandler) getTemplate(ctxt *context, w http.ResponseWriter, req *http.Request) {
	// Check that the template exists
	if h.templates.Lookup(ctxt.Page) == nil {
		http.NotFound(w, req)
		return
	}
	err := h.templates.ExecuteTemplate(w, ctxt.Page, ctxt)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *templateHandler) createDocument(w http.ResponseWriter, req *http.Request) {
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
	http.Redirect(w, req, fmt.Sprintf("/all/%d", doc.Id()), 303)
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
	desc := req.PostFormValue("description")
	tags := strings.Split(req.PostFormValue("tags"), ",")
	err := ctxt.Project.Update(func(pf *project.ProjectFields) {
		pf.Description = desc
		pf.Tags = tags
	})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	// And return the updated document
	err = h.templates.ExecuteTemplate(w, "project-li", []any{ctxt.Project, ctxt})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
}

func (h *templateHandler) updateDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
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
	err = ctxt.Document.Update(func(df *document.DocumentFeilds) {
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
	err = h.templates.ExecuteTemplate(w, "document", ctxt.Document)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
}

func (h *templateHandler) deleteDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
	err := h.DocumentStore.Remove(ctxt.Document)
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
	// And grab the associated function
	switch c.Page {
	case "projects":
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
		sort := "byCreated"
		if sortFunc, ok := c.DocumentSortFunctions[sort]; ok {
			c.DocumentStore = c.DocumentStore.Sort(sortFunc.Func).(*document.Store)
		}
	}
}
