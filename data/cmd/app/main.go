package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pkab00/dockloud/data"
	"github.com/pkab00/dockloud/data/internal/cache"
	"github.com/pkab00/dockloud/data/internal/database"
	"github.com/pkab00/dockloud/shared/middleware"
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

	mux := http.NewServeMux()
	data.RegisterRoutes(db, rd, mux)
	http.ListenAndServe(":8081", middleware.Logger(mux))
}
