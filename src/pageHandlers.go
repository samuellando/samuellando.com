package main

import (
      "./page"
      "./session"
      "net/http"
      "log"
      "strings"
)

const PAGES_DIR = "pages"

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  if p.Load() != nil {
    http.NotFound(w, r)
    return
  }
  u, _ := session.Active(r)
  var uName = "nil"
  if u != nil {
    uName = (*u).UserName()
  }
  if Allowed(uName, r.URL.Path) || Allowed("nil", r.URL.Path) {
    renderTemplate(w, "view", p)
  } else {
    log.Print("User not allowed")
    http.NotFound(w, r)
  }
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  err := p.Load()
  u, _ := session.Active(r)
  var uName string
  if u == nil {
    log.Print("Must be logged in to edit")
    http.NotFound(w, r)
    return
  }
  uName = (*u).UserName()
  if err != nil || Allowed(uName, r.URL.Path) {
    renderTemplate(w, "edit", p)
  } else {
    log.Print("User not allowed")
    http.NotFound(w, r)
  }
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  err := p.Load()
  u, _ := session.Active(r)
  if u == nil {
    log.Print("Must be logged in to save")
    http.NotFound(w, r)
    return
  }
  uName := (*u).UserName()
  oldTitle := strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/")) - 1]
  if err != nil || Allowed(uName, "/save/"+oldTitle) {
    p.Remove()
    title := r.FormValue("title")
    body := r.FormValue("body")
    p = page.New(PAGES_DIR, title, []byte(body))
    DisAllow(uName, "/view/"+oldTitle)
    DisAllow(uName, "/edit/"+oldTitle)
    DisAllow(uName, "/save/"+oldTitle)
    DisAllow("nil", "/view/"+oldTitle)
    Allow(uName, "/edit/"+title)
    Allow(uName, "/save/"+title)
    private := r.FormValue("private")
    if private == "yes" {
      Allow(uName, "/view/"+title)
    } else {
      Allow("nil", "/view/"+title)
    }
    p.Save()
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
  } else {
    log.Print("User allowed")
    http.NotFound(w, r)
  }
}
