package main

import (
      "./page"
      "./session"
      "net/http"
      "log"
)

const PAGES_DIR = "pages"

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  if p.Load() != nil {
    http.NotFound(w, r)
    return
  }
  u, _ := session.Active(r)
  if p.WhiteListed(u) {
    renderTemplate(w, "view", p)
  } else {
    log.Print("User not whitelisted")
    http.NotFound(w, r)
  }
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  p.Load()
  u, _ := session.Active(r)
  if p == nil || p.WhiteListed(u) {
    renderTemplate(w, "edit", p)
  } else {
    log.Print("User not whitelisted")
    http.NotFound(w, r)
  }
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  p := page.New(PAGES_DIR, title)
  p.Load()
  u, _ := session.Active(r)
  if p == nil || p.WhiteListed(u) {
    p.Remove()
    title := r.FormValue("title")
    body := r.FormValue("body")
    p = page.New(PAGES_DIR, title, []byte(body))
    private := r.FormValue("private")
    if private == "yes" {
      p.AddUser(u)
    }
    p.Save()
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
  } else {
    log.Print("User not whitelisted")
    http.NotFound(w, r)
  }

}
