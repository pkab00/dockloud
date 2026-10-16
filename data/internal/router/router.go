package router

import (
	"database/sql"
	"net/http"

	"github.com/pkab00/dockloud/data/internal/accesslink"
	"github.com/pkab00/dockloud/data/internal/document"
	"github.com/pkab00/dockloud/data/internal/serviceauth"
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

	mux.Handle(
		"GET /users/find",
		serviceauth.RequirePermission(serviceauth.ReadUsers)(http.HandlerFunc(userHand.Find)),
	)
	mux.Handle(
		"GET /users/auth",
		serviceauth.RequirePermission(serviceauth.ReadUsers)(http.HandlerFunc(userHand.Auth)),
	)
	mux.Handle(
		"POST /users",
		serviceauth.RequirePermission(serviceauth.CreateUsers)(http.HandlerFunc(userHand.Create)),
	)
	mux.Handle(
		"DELETE /users/{uid}",
		serviceauth.RequirePermission(serviceauth.DeleteUsers)(http.HandlerFunc(userHand.Delete)),
	)
	mux.Handle(
		"PATCH /users/{uid}",
		serviceauth.RequirePermission(serviceauth.UpdateUsers)(http.HandlerFunc(userHand.Update)),
	)

	mux.Handle(
		"GET /documents",
		serviceauth.RequirePermission(serviceauth.ReadDocuments)(http.HandlerFunc(docHand.Get)),
	)
	mux.Handle(
		"GET /documents/{uid}",
		serviceauth.RequirePermission(serviceauth.ReadDocuments)(http.HandlerFunc(docHand.GetByID)),
	)
	mux.Handle(
		"POST /documents",
		serviceauth.RequirePermission(serviceauth.CreateDocuments)(http.HandlerFunc(docHand.Create)),
	)
	mux.Handle(
		"DELETE /documents/{id}",
		serviceauth.RequirePermission(serviceauth.DeleteDocuments)(http.HandlerFunc(docHand.Delete)),
	)
	mux.Handle(
		"PATCH /documents/{id}",
		serviceauth.RequirePermission(serviceauth.UpdateDocuments)(http.HandlerFunc(docHand.Update)),
	)

	mux.Handle(
		"GET /links",
		serviceauth.RequirePermission(serviceauth.ReadLinks)(http.HandlerFunc(linkHand.Get)),
	)
	mux.Handle(
		"GET /links/{token}",
		serviceauth.RequirePermission(serviceauth.ReadLinks)(http.HandlerFunc(linkHand.GetByToken)),
	)
	mux.Handle(
		"POST /links",
		serviceauth.RequirePermission(serviceauth.CreateLinks)(http.HandlerFunc(linkHand.Create)),
	)
	mux.Handle(
		"DELETE /links/{token}",
		serviceauth.RequirePermission(serviceauth.RevokeLinks)(http.HandlerFunc(linkHand.Revoke)),
	)

	// for healthcheck
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
