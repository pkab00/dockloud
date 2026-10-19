package main

import (
	"net/http"

	"github.com/pkab00/dockloud/link/internal/router"
)

func main() {
	mux := http.NewServeMux()

	router.RegisterRoutes(mux)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8082", mux)
}
