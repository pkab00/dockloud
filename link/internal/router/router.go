package router

import (
	"net/http"

	"github.com/pkab00/dockloud/link/internal/client/data"
	"github.com/pkab00/dockloud/link/internal/document"
	"github.com/pkab00/dockloud/link/internal/link"
)

func RegisterRoutes(mux *http.ServeMux) {
	dataClient, err := data.NewClient("http://data:8081")
	if err != nil {
		panic(err)
	}

	var (
		linkService = link.NewService(dataClient)
		linkHand    = link.NewHandler(linkService)

		docService = document.NewService(dataClient)
		docHand    = document.NewHandler(docService)
	)

	mux.HandleFunc("GET /links", linkHand.Get)
	mux.HandleFunc("POST /links", linkHand.Create)
	mux.HandleFunc("DELETE /links/{token}", linkHand.Revoke)

	mux.HandleFunc("GET /documents/{token}", docHand.Get)
	mux.HandleFunc("PATCH /documents/{token}", docHand.Update)
	mux.HandleFunc("DELETE /documents/{token}", docHand.Delete)
}
