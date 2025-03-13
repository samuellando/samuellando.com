package search

import (
	"html/template"
	"net/http"
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
