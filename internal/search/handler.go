package search

import (
    "net/http"
    "html/template"
    "strings"
)

func CreateSearchHandler(indexes ...func() []IndexItem) http.HandlerFunc {
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
		if searchString == "" {
			return
		}
		for _, index := range indexes {
			elements := index()
			for _, elem := range elements {
				if strings.Contains(elem.Text, searchString) {
                    t.Execute(w, elem)
				}
			}
		}
	}
}
