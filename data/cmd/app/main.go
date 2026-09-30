package main

import (
	"log"
	"os"

	"github.com/pkab00/dockloud/data/internal/cache"
	"github.com/pkab00/dockloud/data/internal/database"
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
}
