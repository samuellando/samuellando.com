package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "html/template"
    "regexp"
    )

type Page struct {
  Title string
  Body []byte
}

var templatesDir = "tmpl"
var pagesDir = "pages"

var validPath = regexp.MustCompile("^/(edit|save|view|static)/([a-zA-Z0-9.]+)$")

var templates = template.Must(template.ParseFiles(templatesDir + "/edit.html",
                                                  templatesDir + "/view.html",
                                                  templatesDir + "/home.html",
                                                  templatesDir + "/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func (p *Page) savePage() error {
  filename := pagesDir + "/" + p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := pagesDir + "/" + title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}

func viewHandler (w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    fmt.Fprint(w, "Error loading page")
  } else {
    renderTemplate(w, "view" ,p)
  }
}

func editHandler (w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit" ,p)
}

func saveHandler (w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  p.savePage()
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func staticHandler (w http.ResponseWriter, r *http.Request, file string) {
  http.ServeFile(w, r, "static/"+file)
}

func indexHandler (w http.ResponseWriter, r *http.Request) {
  files, err := ioutil.ReadDir(pagesDir)
  if err != nil {
   log.Fatal(err)
  }
  var names = make([]string, 0)
  for i := 0; i < len(files); i++ {
    names = append(names, files[i].Name()[:len(files[i].Name())-4])
  }
  err = templates.ExecuteTemplate(w, "index.html", names)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func homeHandler (w http.ResponseWriter, r *http.Request) {
  renderTemplate(w, "home", nil)
}

func main() {
  http.HandleFunc("/", homeHandler)
  http.HandleFunc("/index", indexHandler)
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.HandleFunc("/static/", makeHandler(staticHandler))
  log.Fatal(http.ListenAndServe(":8080", nil))
}
