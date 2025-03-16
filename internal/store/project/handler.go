package project

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"html/template"
	"samuellando.com/internal/store/tag"
)

type Handler struct {
	Template     template.Template
	ProjectStore Store
	TagStore     tag.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.templateRequest(w, req)
	case "PUT":
		h.updateProject(w, req)
	}
}

func (h *Handler) templateRequest(w http.ResponseWriter, req *http.Request) {
	proj := h.getReqProject(req)
	h.renderProject(w, proj)
}

func (h *Handler) renderProject(w http.ResponseWriter, project Project) {
	err := h.Template.Execute(w, project)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

func (h *Handler) getReqProject(req *http.Request) Project {
	id, err := strconv.Atoi(req.PathValue("project"))
	if err != nil {
		return Project{}
	}
	proj, err := h.ProjectStore.GetById(int64(id))
	if err != nil {
		return Project{}
	}
	return proj
}

func (h *Handler) getTagsFromReq(req *http.Request) []tag.ProtoTag {
	tagValues := strings.Split(req.PostFormValue("tags"), ",")
	tags := make([]tag.ProtoTag, 0)
	for _, tv := range tagValues {
		if tv == "" {
			continue
		}
		tags = append(tags, tag.ProtoTag{Value: tv})
	}
	return tags
}

func (h *Handler) updateProject(w http.ResponseWriter, req *http.Request) {
	proj := h.getReqProject(req)
	rdesc := req.PostFormValue("description")
	rimage := req.PostFormValue("image")
	rhidden := req.PostFormValue("hidden")
	tags := h.getTagsFromReq(req)
	var desc *string
	if rdesc != "" {
		desc = &rdesc
	}
	var image *string
	if rimage != "" {
		image = &rimage
	}
	hidden := false
	if rhidden == "true" {
		hidden = true
	}
	err := proj.Update(func(pf *ProtoProject) {
		pf.Description = desc
		pf.ImageLink = image
		pf.Tags = tags
		pf.Hidden = hidden
	})
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
		return
	}
	h.renderProject(w, proj)
}
