package asset

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	Store Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.getAsset(w, req)
	case "POST":
		h.createAsset(w, req)
	case "DELETE":
		h.deleteAsset(w, req)
	}
}

func (h *Handler) getAsset(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("asset")
	if name == "" {
		http.Error(w, "asset name must be provided", 404)
	} else {
		asset, err := h.Store.GetByName(name)
		if err != nil {
			http.Error(w, "asset not found", 404)
			return
		}
		content, err := asset.Content()
		if err != nil {
			http.Error(w, "Unable to load content", 500)
			return
		}
		w.Write(content)
	}
}

func (h *Handler) createAsset(w http.ResponseWriter, req *http.Request) {
	const max_file_size = int64(4000000)
	f, header, err := req.FormFile("file")
	// If there is no file
	if err != nil {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(400), "No file provided"), 400)
		return
	}
	if header.Size > max_file_size {
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(413), "File too large (4MB max)"), 413)
		return
	}
	defer f.Close()
	buff := make([]byte, header.Size)
	for {
		r := bufio.NewReader(f)
		_, err = r.Read(buff)
		if err != nil && err != io.EOF {
			http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
			return
		}
		if err != io.EOF {
			break
		}
	}
	_, err = h.Store.Add(ProtoAsset{
		Name:    header.Filename,
		Content: buff,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
}

func (h *Handler) deleteAsset(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("asset")
	asset, err := h.Store.GetByName(name)
	if err != nil {
		http.Error(w, "Could not find asset", 404)
		return
	}
	err = asset.Delete()
	if err != nil {
		http.Error(w, "Faild to delete asset", 500)
		return
	}
}
