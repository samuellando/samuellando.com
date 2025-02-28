package search

import (
	"html/template"
	"net/http"
)

func CreateSearchHandler(indexes ...indexFunc) http.HandlerFunc {
	tmpl := `
    <div>
    {{.Type}} <a href="{{.Path}}">{{.Item.Title}}</a>
    </div>
    `
	t, err := template.New("result").Parse(tmpl)
	if err != nil {
		panic(err)
	}
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
