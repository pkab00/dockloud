package router

import (
	"log"
	"net/http"

	"github.com/pkab00/dockloud/api/internal/auth"
	"github.com/pkab00/dockloud/api/internal/client/data"
	"github.com/pkab00/dockloud/api/internal/document"
	"github.com/pkab00/dockloud/api/internal/middleware"
)

func RegisterRoutes(mux *http.ServeMux, JWTsecret string) {
	var err error

	dataClient, err := data.NewClient("http://data:8081")
	if err != nil {
		log.Fatalf("failed to initialize data client: %v", err)
	}
	verifier, err := auth.NewVerifier(JWTsecret)
	if err != nil {
		log.Fatalf("failed to initialize verifier: %v", err)
	}

	var (
		authService = auth.NewService(dataClient)
		authHandler = auth.NewHandler(authService)
		docService  = document.NewService(dataClient)
		docHandler  = document.NewHandler(docService)
	)

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle(
		"POST /documents",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.Create)),
	)
	mux.Handle(
		"GET /documents",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.GetAll)),
	)
}
