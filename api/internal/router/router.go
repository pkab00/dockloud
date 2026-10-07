package router

import (
	"net/http"

	"github.com/pkab00/dockloud/api/internal/auth"
	"github.com/pkab00/dockloud/api/internal/client/data"
)

func RegisterRoutes(mux *http.ServeMux) {
	var (
		dataClient = data.NewClient("http://data:8081")

		authService = auth.NewService(dataClient)
		authHandler = auth.NewHandler(authService)
	)

	mux.HandleFunc("POST /auth", authHandler.CreateUser)
}
