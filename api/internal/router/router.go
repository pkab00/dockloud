package router

import (
	"net/http"

	"github.com/pkab00/dockloud/api/internal/auth"
	"github.com/pkab00/dockloud/api/internal/client/data"
)

func RegisterRoutes(mux *http.ServeMux) {
	var dataClient, err = data.NewClient("http://data:8081")
	if err != nil {
		panic("Failed to initialize data client")
	}

	var (
		authService = auth.NewService(dataClient)
		authHandler = auth.NewHandler(authService)
	)

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
}
