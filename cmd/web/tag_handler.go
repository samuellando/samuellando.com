package main

import (
	"fmt"
	"net/http"
	"strconv"

	"samuellando.com/internal/store/tag"
	"samuellando.com/internal/template"
)

type tagHandler struct {
	Store     *tag.Store
	Templates template.Template
}

func (h *tagHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PATCH":
		h.updateTag(w, req)
	case "DELETE":
		h.deleteTag(w, req)
	}
}

func (h *tagHandler) updateTag(w http.ResponseWriter, req *http.Request) {
	ids := req.PathValue("tag")
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
	t, err := h.Store.GetById(id)
	if err != nil {
		http.Error(w, "Could not find tag", 404)
		return
	}
	color := req.FormValue("color")
	err = t.Update(func(tf *tag.TagFields) {
		tf.Color = color
	})
	if err != nil {
		http.Error(w, "Faild to update tag", 500)
		return
	}
}

func (h *tagHandler) deleteTag(w http.ResponseWriter, req *http.Request) {
	ids := req.PathValue("tag")
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
	asset, err := h.Store.GetById(id)
	if err != nil {
		http.Error(w, "Could not find tag", 404)
		return
	}
	err = asset.Delete()
	if err != nil {
		http.Error(w, "Faild to delete tag", 500)
		return
	}
}
