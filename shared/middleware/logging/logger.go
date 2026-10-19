package logging

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(lw, r)
		log.Printf("%s %s | %d", r.Method, r.RequestURI, lw.status)
	})
}
