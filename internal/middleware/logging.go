package middleware

import (
	"log"
	"net/http"
	"time"
)

type logResponseWriter struct {
	writer http.ResponseWriter
	code   int
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *logResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.writer.WriteHeader(statusCode)
}

func LoggingFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logWriter := &logResponseWriter{writer: w, code: 200}
		log.Printf("Request: [%s] %s\n", req.Method, req.URL.Path)
		start := time.Now()
		h(logWriter, req)
		log.Printf("Response: %d (in %s)", logWriter.code, time.Since(start))
	}
}

func Logging(h http.Handler) http.HandlerFunc {
	return LoggingFunc(h.ServeHTTP)
}
