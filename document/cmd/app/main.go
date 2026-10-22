package main

import (
	"net/http"
	"os"

	"github.com/pkab00/dockloud/document/internal/router"
	"github.com/pkab00/dockloud/shared/middleware/logging"
)

func main() {
	JWTsecret := os.Getenv("JWT_KEY")
	mux := http.NewServeMux()

	router.RegisterRoutes(mux, JWTsecret)
	http.ListenAndServe(":8080", logging.Logger(mux))
}
