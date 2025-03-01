package search

import (
	"net/http"
	"samuellando.com/internal/template"
)

func CreateSearchHandler(t template.Template, indexes ...indexFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		searchString := req.FormValue("q")
		for i, elem := range searchIndexes(searchString, indexes...) {
			if i > 2 {
				break
			}
			t.Execute(w, elem)
		}
	}
}
