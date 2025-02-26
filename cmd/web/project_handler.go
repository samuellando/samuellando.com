package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"samuellando.com/internal/store/project"
	"samuellando.com/internal/template"
)

type projectHandler struct {
	templates    template.Template
	projectStore project.Store
}

func (h *projectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.templateRequest(w, req)
	case "PUT":
		h.updateProject(w, req)
	}
}

func (h *projectHandler) templateRequest(w http.ResponseWriter, req *http.Request) {
	doc := h.getReqProject(req)
	h.renderProject(w, doc)
}

func (h *projectHandler) renderProject(w http.ResponseWriter, project *project.Project) {
	err := h.templates.ExecuteTemplate(w, "project", project)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *projectHandler) getReqProject(req *http.Request) *project.Project {
	id, err := strconv.Atoi(req.PathValue("project"))
	if err != nil {
		return nil
	}
	doc, err := h.projectStore.GetById(id)
	if err != nil {
		return nil
	}
	return doc
}

func (h *projectHandler) updateProject(w http.ResponseWriter, req *http.Request) {
	proj := h.getReqProject(req)
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
	h.renderProject(w, proj)
}
