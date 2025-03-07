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
	proj := h.getReqProject(req)
	h.renderProject(w, proj)
}

func (h *projectHandler) renderProject(w http.ResponseWriter, project project.Project) {
	err := h.templates.ExecuteTemplate(w, "project", project)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *projectHandler) getReqProject(req *http.Request) project.Project {
	id, err := strconv.Atoi(req.PathValue("project"))
	if err != nil {
		return project.Project{}
	}
	proj, err := h.projectStore.GetById(int64(id))
	if err != nil {
		return project.Project{}
	}
	return proj
}

func (h *projectHandler) getTagsFromReq(req *http.Request) []tag.ProtoTag {
	tagValues := strings.Split(req.PostFormValue("tags"), ",")
	tags := make([]tag.ProtoTag, len(tagValues))
	for i, tv := range tagValues {
		tags[i] = tag.ProtoTag{
			Value: tv,
		}
	}
	return tags
}

func (h *projectHandler) updateProject(w http.ResponseWriter, req *http.Request) {
	proj := h.getReqProject(req)
	rdesc := req.PostFormValue("description")
	tags := h.getTagsFromReq(req)
	var desc *string
	if rdesc != "" {
		desc = &rdesc
	}
	err := proj.Update(func(pf *project.ProtoProject) {
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
