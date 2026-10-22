package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pkab00/dockloud/data/internal/cache"
	"github.com/pkab00/dockloud/data/internal/database"
	"github.com/pkab00/dockloud/data/internal/router"
	"github.com/pkab00/dockloud/shared/middleware/logging"
	"github.com/pkab00/dockloud/shared/middleware/service_auth"
)

func main() {
	db, err := database.NewDB(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rd, err := cache.NewCache(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer rd.Close()

	verifier, err := serviceauth.NewVerifier(
		map[string]serviceauth.ServiceConfig{
			"dockloud-api": {
				Secret:   []byte(os.Getenv("API_SERVICE_KEY")),
				Subject:  "dockloud-api",
				Audience: "dockloud-data",
			},
			"dockloud-link": {
				Secret:   []byte(os.Getenv("LINK_SERVICE_KEY")),
				Subject:  "dockloud-link",
				Audience: "dockloud-data",
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	router.RegisterRoutes(db, rd, mux)
	http.ListenAndServe(
		":8081",
		logging.Logger(
			serviceauth.Middleware(verifier)(mux),
		),
	)
}
