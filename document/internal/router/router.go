package router

import (
	"log"
	"net/http"

	"github.com/pkab00/dockloud/document/internal/auth"
	dataclient "github.com/pkab00/dockloud/document/internal/client/data"
	linkclient "github.com/pkab00/dockloud/document/internal/client/link"
	"github.com/pkab00/dockloud/document/internal/document"
	"github.com/pkab00/dockloud/document/internal/link"
	"github.com/pkab00/dockloud/document/internal/middleware"
)

func RegisterRoutes(mux *http.ServeMux, JWTsecret string) {
	var err error

	dataClient, err := dataclient.NewClient("http://data:8081")
	if err != nil {
		log.Fatalf("failed to initialize data client: %v", err)
	}
	linkClient, err := linkclient.NewClient("http://link:8082")
	if err != nil {
		log.Fatalf("failed to initialize link client: %v", err)
	}
	verifier, err := auth.NewVerifier(JWTsecret)
	if err != nil {
		log.Fatalf("failed to initialize verifier: %v", err)
	}

	var (
		authService = auth.NewService(dataClient)
		authHandler = auth.NewHandler(authService)
		docService  = document.NewService(dataClient, linkClient)
		docHandler  = document.NewHandler(docService)
		linkService = link.NewService(linkClient)
		linkHandler = link.NewHandler(linkService)
	)

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	mux.Handle(
		"POST /documents",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.Create)),
	)
	mux.Handle(
		"GET /documents",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.Get)),
	)
	mux.Handle(
		"GET /documents/{id}",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.GetByID)),
	)
	mux.Handle(
		"PATCH /documents/{id}",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.Update)),
	)
	mux.Handle(
		"DELETE /documents/{id}",
		middleware.Auth(verifier)(http.HandlerFunc(docHandler.Delete)),
	)

	mux.Handle(
		"GET /documents/shared/{token}",
		http.HandlerFunc(docHandler.GetShared),
	)
	mux.Handle(
		"PATCH /documents/shared/{token}",
		http.HandlerFunc(docHandler.UpdateShared),
	)
	mux.Handle(
		"DELETE /documents/shared/{token}",
		http.HandlerFunc(docHandler.DeleteShared),
	)

	mux.Handle(
		"GET /links",
		middleware.Auth(verifier)(http.HandlerFunc(linkHandler.Get)),
	)
	mux.Handle(
		"POST /links",
		middleware.Auth(verifier)(http.HandlerFunc(linkHandler.Create)),
	)
	mux.Handle(
		"DELETE /links/{token}",
		middleware.Auth(verifier)(http.HandlerFunc(linkHandler.Revoke)),
	)
}
