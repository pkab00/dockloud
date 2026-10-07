package main

import (
	"net/http"

	"github.com/pkab00/dockloud/api/internal/router"
	"github.com/pkab00/dockloud/shared/middleware"
)

func main() {
	mux := http.NewServeMux()
	router.RegisterRoutes(mux)
	http.ListenAndServe(":8080", middleware.Logger(mux))
}
