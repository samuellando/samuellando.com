package main

import (
	"fmt"
	"net/http"
	"samuellando.com/internal/auth"
)

func authenticate(w http.ResponseWriter, req *http.Request) {
	if auth.ValidCredentials(req) {
		cookie := &http.Cookie{
			Name:  "session",
			Value: auth.CreateJWT(),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, req, "/", 303)
		return
	} else {
		fmt.Fprint(w, "Invalid login dredentials")
		http.Redirect(w, req, "/", 401)
		return
	}
}

func deauthenticate(w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", 303)
}

