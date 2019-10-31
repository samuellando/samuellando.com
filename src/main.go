package main

import (
	"./page"
	"./session"
	"./user"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

const USERS_DB = "users.db"

var validUrl = regexp.MustCompile("^/(edit|save|view|static)/([a-zA-Z0-9.]+)$")

const TEMPLATES_DIR = "tmpl"

var templates = template.Must(template.ParseFiles(
	TEMPLATES_DIR+"/titlebar.html",
	TEMPLATES_DIR+"/view.html",
	TEMPLATES_DIR+"/edit.html",
	TEMPLATES_DIR+"/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
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
        u, _ := session.Active(r)
        renderTemplate(w, "titlebar", u)
	pages := page.List(PAGES_DIR, u)
	renderTemplate(w, "index", pages)
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
		http.Redirect(w, r, "/static/login.html", http.StatusFound)
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
		http.Redirect(w, r, "/static/signup.html", http.StatusFound)
	}
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {
	session.Destroy(w, r)
	http.Redirect(w, r, "/index", http.StatusFound)
}

func main() {
	http.HandleFunc("/static/", makeHandler(staticHandler))
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/login", logInHandler)
	http.HandleFunc("/logout", logOutHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
