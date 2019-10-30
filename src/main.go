package main

import (
	"./page"
	"./session"
	"./user"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

const USERS_DB = "users.db"

var validUrl = regexp.MustCompile("^/(edit|save|view|static)/([a-zA-Z0-9.]+)$")

const TEMPLATES_DIR = "tmpl"

var templates = template.Must(template.ParseFiles(
	TEMPLATES_DIR+"/view.html",
	TEMPLATES_DIR+"/save.html",
	TEMPLATES_DIR+"/edit.html",
	TEMPLATES_DIR+"/home.html",
	TEMPLATES_DIR+"/index.html",
	TEMPLATES_DIR+"/login.html",
	TEMPLATES_DIR+"/signup.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p page.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validUrl.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request, file string) {
	http.ServeFile(w, r, "static/"+file)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	pages := page.List(PAGES_DIR)
	err := templates.ExecuteTemplate(w, "index.html", pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	u, err := session.Active(r)
	if err != nil {
		fmt.Fprintf(w, "Not logged in")
	} else {
		fmt.Fprintf(w, "Logged in as: %s", (*u).UserName())
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", nil)
}

func logInHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	if userName != "" && password != "" {
		u := user.New(USERS_DB, userName)
		err := u.Validate(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Create(w, r, u)
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		renderTemplate(w, "login", nil)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	if userName != "" && password != "" {
		u := user.New(USERS_DB, userName)
		err := u.Add(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		session.Create(w, r, u)
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		renderTemplate(w, "signup", nil)
	}
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {
	session.Destroy(w, r)
	http.Redirect(w, r, "/index", http.StatusFound)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/login", logInHandler)
	http.HandleFunc("/logout", logOutHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/static/", makeHandler(staticHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
