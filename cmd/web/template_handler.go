package main

import (
	"bufio"
	"fmt"
    "io"
	"html/template"
	"log"
	"net/http"
	"samuellando.com/internal/auth"
	"samuellando.com/internal/stores"
	"strconv"
	"strings"
)

type context struct {
    Handler     *templateHandler
	Page        string
	DocumentRef string
	Document    stores.Document
	Admin       bool
}

func getPathContext(req *http.Request) *context {
	path := req.URL.Path
	var page string
	var documentRef string
	if path == "/" {
		page = "index"
	} else {
		page = req.PathValue("template")
		documentRef = req.PathValue("document")
	}
	return &context{Page: page, DocumentRef: documentRef, Admin: isAuthenticated(req)}
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
	MarkdownStore stores.MarkdownStore
	ProjectStore  stores.ProjectStore
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctxt := getPathContext(req)
    ctxt.Handler = h
	if ctxt.DocumentRef != "" {
		documentId, err := strconv.Atoi(ctxt.DocumentRef)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(400), "document reference must be numeric"), 400)
			return
		}
		ctxt.Document, err = h.MarkdownStore.GetDocumentById(documentId)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(404), "document not found"), 404)
			return
		}
	}
	switch req.Method {
	case "GET":
		h.getTemplate(ctxt, w, req)
	case "POST":
		h.createDocument(ctxt, w, req)
	case "PUT":
		h.updateDocument(ctxt, w, req)
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

func (h *templateHandler) createDocument(context *context, w http.ResponseWriter, req *http.Request) {
	title := req.PostFormValue("title")
	content := req.PostFormValue("content")
	tags := strings.Split(req.PostFormValue("tags"), ",")
	document, err := h.MarkdownStore.CreateDocument(title, content, tags)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	http.Redirect(w, req, fmt.Sprintf("/all/%d", document.Id()), 303)
}

func (h *templateHandler) updateDocument(ctxt *context, w http.ResponseWriter, req *http.Request) {
	max_file_size := int64(2000000)
	if f, header, err := req.FormFile("file"); err == nil {
		defer f.Close()
		if header.Size > max_file_size {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(400), "File too large (2MB max)"), 400)
			return
		}
		buff := make([]byte, header.Size)
		for {
			r := bufio.NewReader(f)
			_, err = r.Read(buff)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if err != io.EOF {
				break
			}
		}
		err = ctxt.Document.SetContent(string(buff))
		w.Header().Add("HX-Refresh", "true")
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
	} else {
		err = ctxt.Document.SetContent(req.PostFormValue("content"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
	}
	err := ctxt.Document.SetTitle(req.PostFormValue("title"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	if req.PostFormValue("published") == "true" {
		err = ctxt.Document.SetPublished(true)
	} else {
		err = ctxt.Document.SetPublished(false)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	tags := strings.Split(req.PostFormValue("tags"), ",")
	err = ctxt.Document.SetTags(tags)
	if err != nil {
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
	err := ctxt.Document.Delete()
	// Failed to delete
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	http.Redirect(w, req, "/"+ctxt.Page, 303)
	return
}

