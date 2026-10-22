package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pkab00/dockloud/link/internal/router"
	"github.com/pkab00/dockloud/shared/middleware/logging"
	serviceauth "github.com/pkab00/dockloud/shared/middleware/service_auth"
)

func main() {
	mux := http.NewServeMux()

	router.RegisterRoutes(mux)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	verifier, err := serviceauth.NewVerifier(
		map[string]serviceauth.ServiceConfig{
			"dockloud-api": {
				Secret:   []byte(os.Getenv("API_SERVICE_KEY")),
				Subject:  "dockloud-api",
				Audience: "dockloud-link",
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(
		":8082",
		logging.Logger(serviceauth.Middleware(verifier)(mux)),
	)
}
