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
    "samuellando.com/internal/stores"
    "samuellando.com/internal/auth"
)

type handler struct {
	templates     template.Template
	MarkdownStore stores.MarkdownStore
	ProjectStore  stores.ProjectStore
	assetsServer  http.Handler
}

type context struct {
	Handler  *handler
	Request  *http.Request
	Page     string
	Document stores.Document
	Admin    bool
}

func createHandler(templateDir, assetsDir, assetsPrefix string) *handler {
	templates, err := template.New("templates").Funcs(template.FuncMap{"join": strings.Join}).ParseGlob(templateDir + "/*")
	if err != nil {
		panic(err)
	}
	ms := stores.InitializeMarkdownStore()
	ps := stores.InitializeProjectStore()
	assetsServer := http.StripPrefix(assetsPrefix, http.FileServer(http.Dir(assetsDir)))
	return &handler{MarkdownStore: ms, templates: *templates, assetsServer: assetsServer, ProjectStore: ps}
}

func (h *handler) RenderMarkdown(page string) template.HTML {
	return template.HTML("<h1>Hello!</h1><p>My name is sam</p>" + "PAGE=" + page)
}

func isAuthenticated(req *http.Request) bool {
	if cookie, err := req.Cookie("session"); err == nil && auth.ValidJWT(cookie.Value) {
		return true
	} else {
		return false
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	if strings.HasPrefix(path, "/assets/") && req.Method == "GET" {
		h.assetsServer.ServeHTTP(w, req)
		return
	}
	if path == "/auth" && req.Method == "POST" {
		// TODO : Set the tright values!
		if auth.ValidCredentials(req) {
			cookie := &http.Cookie{
				Name:  "session",
				Value: auth.CreateJWT(),
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, req, "/", 303)
			return
		} else {
			fmt.Fprint(w, "Invalid login dredentials")
			http.Redirect(w, req, "/", 401)
			return
		}
	}
	if path == "/deauth" && req.Method == "POST" {
		cookie := &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, req, "/", 303)
		return
	}
	log.Println(method, path)
	var page string
	var document stores.Document
	if path == "/" {
		page = "index"
	} else {
		parts := strings.Split(path, "/")
		page = parts[1]
		if len(parts) >= 3 {
			var err error
			documentId, err := strconv.Atoi(parts[2])
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(400), "Document id must be numeric"), 400)
				return
			}
			document, err = h.MarkdownStore.GetDocumentById(documentId)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(404), err), 404)
				return
			}
		}
	}
	context := context{h, req, page, document, isAuthenticated(req)}
	switch req.Method {
	case "GET":
		// Check that the template exists
		if h.templates.Lookup(page) == nil {
			http.NotFound(w, req)
			return
		}
		err := h.templates.ExecuteTemplate(w, page, context)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
	case "DELETE":
		err := document.Delete()
		// Failed to delete
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
		http.Redirect(w, req, "/"+page, 303)
		return
	case "PUT":
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
			err = document.SetContent(string(buff))
			w.Header().Add("HX-Refresh", "true")
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
				return
			}
		} else {
			err = document.SetContent(req.PostFormValue("content"))
			if err != nil {
				http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
				return
			}
		}
		err := document.SetTitle(req.PostFormValue("title"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
		if req.PostFormValue("published") == "true" {
			err = document.SetPublished(true)
		} else {
			err = document.SetPublished(false)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
		tags := strings.Split(req.PostFormValue("tags"), ",")
		err = document.SetTags(tags)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
		// And return the updated document
		err = h.templates.ExecuteTemplate(w, "document", document)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
	case "POST":
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
}

func (h *handler) close() {
	h.MarkdownStore.Close()
}

func main() {
	handler := createHandler("./templates", "./assets", "/assets")
	defer handler.close()
	http.ListenAndServe(":8080", handler)
}
