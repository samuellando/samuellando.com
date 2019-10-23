package main

import (
    "log"
    "./page"
    "net/http"
    "html/template"
    "regexp"
    "io/ioutil"
    )

const PAGES_DIR = "pages"

var validUrl = regexp.MustCompile("^/(edit|save|view|static)/([a-zA-Z0-9.]+)$")

const TEMPLATES_DIR = "tmpl"
var templates = template.Must(template.ParseFiles(TEMPLATES_DIR + "/edit.html",
                                                  TEMPLATES_DIR + "/view.html",
                                                  TEMPLATES_DIR + "/home.html",
                                                  TEMPLATES_DIR + "/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p page.Page) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validUrl.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}

func viewHandler (w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  err := p.Load()
  if err != nil {
    log.Print(err)
    http.NotFound(w, r)
    return
  } else {
    renderTemplate(w, "view", p)
  }
}

func editHandler (w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  err := p.Load()
  if err != nil {
    http.Redirect(w, r, "/new/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "edit" ,p)
}

func saveHandler (w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := page.New(PAGES_DIR, title, []byte(body))
  p.Save()
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func staticHandler (w http.ResponseWriter, r *http.Request, file string) {
  http.ServeFile(w, r, "static/"+file)
}

func indexHandler (w http.ResponseWriter, r *http.Request) {
  files, err := ioutil.ReadDir(PAGES_DIR)
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
