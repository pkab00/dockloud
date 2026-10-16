package router

import (
	"database/sql"
	"net/http"

	"github.com/pkab00/dockloud/data/internal/accesslink"
	"github.com/pkab00/dockloud/data/internal/document"
	"github.com/pkab00/dockloud/data/internal/user"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(db *sql.DB, rd *redis.Client, mux *http.ServeMux) {
	var (
		userRepo = user.NewRepositoryWithCache(db, rd)
		userServ = user.NewService(userRepo)
		userHand = user.NewHandler(userServ)

		docRepo = document.NewRepositoryWithCache(db, rd)
		docServ = document.NewService(docRepo, userRepo)
		docHand = document.NewHandler(docServ)

		linkRepo = accesslink.NewRepository(db)
		linkServ = accesslink.NewService(linkRepo)
		linkHand = accesslink.NewHandler(linkServ)
	)

	mux.HandleFunc("GET /users/find", userHand.Find)
	mux.HandleFunc("GET /users/auth", userHand.Auth)
	mux.HandleFunc("POST /users", userHand.Create)
	mux.HandleFunc("DELETE /users/{uid}", userHand.Delete)
	mux.HandleFunc("PATCH /users/{uid}", userHand.Update)

	mux.HandleFunc("GET /documents", docHand.Get)
	mux.HandleFunc("GET /documents/{id}", docHand.GetByID)
	mux.HandleFunc("POST /documents", docHand.Create)
	mux.HandleFunc("DELETE /documents/{id}", docHand.Delete)
	mux.HandleFunc("PATCH /documents/{id}", docHand.Update)

	mux.HandleFunc("GET /links", linkHand.Get)
	mux.HandleFunc("GET /links/{token}", linkHand.GetByToken)
	mux.HandleFunc("POST /links", linkHand.Create)
	mux.HandleFunc("DELETE /links/{token}", linkHand.Revoke)

	// for healthcheck
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
