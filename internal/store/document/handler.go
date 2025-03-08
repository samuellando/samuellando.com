package document

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"html/template"
	"samuellando.com/internal/store/tag"
)

type Handler struct {
	Template      template.Template
	DocumentStore Store
	TagStore      tag.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if req.FormValue("download") != "" {
			h.downloadDocument(w, req)
		} else {
			h.templateRequest(w, req)
		}
	case "POST":
		h.createDocument(w, req)
	case "PUT":
		h.updateDocument(w, req)
	case "DELETE":
		h.deleteDocument(w, req)
	}
}

func (h *Handler) templateRequest(w http.ResponseWriter, req *http.Request) {
	doc := h.getReqDoc(req)
	h.renderDocument(w, doc)
}

func (h *Handler) renderDocument(w http.ResponseWriter, doc Document) {
	err := h.Template.Execute(w, doc)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *Handler) getReqDoc(req *http.Request) Document {
	id, err := strconv.Atoi(req.PathValue("document"))
	if err != nil {
		return Document{}
	}
	doc, err := h.DocumentStore.GetById(int64(id))
	if err != nil {
		return Document{}
	}
	return doc
}

func (h *Handler) downloadDocument(w http.ResponseWriter, req *http.Request) {
	doc := h.getReqDoc(req)
	filename := fmt.Sprintf("%s.md", doc.Title())
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "text/markdown")
	w.WriteHeader(http.StatusOK)
	content := doc.Content()
	_, err := w.Write([]byte(content))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *Handler) createDocument(w http.ResponseWriter, req *http.Request) {
	title := req.PostFormValue("title")
	content := req.PostFormValue("content")
	tags := h.getTagsFromReq(req)
	doc, err := h.DocumentStore.Add(ProtoDocument{
		Title:   title,
		Content: content,
		Tags:    tags,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	h.renderDocument(w, doc)
}
func (h *Handler) getTagsFromReq(req *http.Request) []tag.ProtoTag {
	tagValues := strings.Split(req.PostFormValue("tags"), ",")
	tags := make([]tag.ProtoTag, len(tagValues))
	for i, tv := range tagValues {
		tags[i] = tag.ProtoTag{
			Value: tv,
		}
	}
	return tags
}

func (h *Handler) updateDocument(w http.ResponseWriter, req *http.Request) {
	doc := h.getReqDoc(req)
	title := req.PostFormValue("title")
	content, err, err_code := getUploadContent(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), err_code)
		return
	}
	if content == "" {
		content = req.PostFormValue("content")
	}
	tags := h.getTagsFromReq(req)
	err = doc.Update(func(df *ProtoDocument) {
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
	h.renderDocument(w, doc)
}

func (h *Handler) deleteDocument(w http.ResponseWriter, req *http.Request) {
	doc := h.getReqDoc(req)
	err := doc.Delete()
	// Failed to delete
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	fmt.Fprint(w, "ok")
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
