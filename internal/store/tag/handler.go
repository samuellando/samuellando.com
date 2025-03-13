package tag

import (
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	Store Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PATCH":
		h.updateTag(w, req)
	case "DELETE":
		h.deleteTag(w, req)
	}
}

func (h *Handler) updateTag(w http.ResponseWriter, req *http.Request) {
	ids := req.PathValue("tag")
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
	t, err := h.Store.GetById(int64(id))
	if err != nil {
		http.Error(w, "Could not find tag", 404)
		return
	}
	color := req.FormValue("color")
	err = t.Update(func(tf *ProtoTag) {
		tf.Color = color
	})
	if err != nil {
		http.Error(w, "Faild to update tag", 500)
		return
	}
}

func (h *Handler) deleteTag(w http.ResponseWriter, req *http.Request) {
	ids := req.PathValue("tag")
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
	asset, err := h.Store.GetById(int64(id))
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
