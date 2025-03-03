package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"samuellando.com/internal/store/project"
	"samuellando.com/internal/store/tag"
	"samuellando.com/internal/template"
)

type projectHandler struct {
	templates    template.Template
	projectStore project.Store
	tagStore     tag.Store
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

func (h *projectHandler) getTagsFromReq(req *http.Request) []tag.Tag {
	tagValues := strings.Split(req.PostFormValue("tags"), ",")
	tags := make([]tag.Tag, len(tagValues))
	for i, tv := range tagValues {
		t, err := h.tagStore.GetByValue(tv)
		if err == nil {
			tags[i] = t
		}
	}
	return tags
}

func (h *projectHandler) updateProject(w http.ResponseWriter, req *http.Request) {
	proj := h.getReqProject(req)
	desc := req.PostFormValue("description")
	tags := h.getTagsFromReq(req)
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
