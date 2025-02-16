package middleware

import (
    "fmt"
	"net/http"
    "log"
	"samuellando.com/internal/auth"
)

func AuthenticatedFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
        if auth.IsAuthenticated(req, w) {
            log.Println("Auth: Validated authentication")
            h(w, req)
        } else {
            http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(403), "Invalid JWT"), 403)
            return
        }
	}
}

func Authenticated(h http.Handler) http.HandlerFunc {
    return AuthenticatedFunc(h.ServeHTTP)
}
