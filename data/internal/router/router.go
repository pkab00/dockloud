package router

import (
	"database/sql"
	"net/http"

	"github.com/pkab00/dockloud/data/internal/accesslink"
	"github.com/pkab00/dockloud/data/internal/document"
	"github.com/pkab00/dockloud/data/internal/permissions"
	"github.com/pkab00/dockloud/data/internal/user"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(db *sql.DB, rd *redis.Client, mux *http.ServeMux) {
	var (
		userRepo = user.NewRepositoryWithCache(db, rd)
		userServ = user.NewService(userRepo)
		userHand = user.NewHandler(userServ)

		docRepo = document.NewRepositoryWithCache(db, rd)
		docServ = document.NewService(docRepo)
		docHand = document.NewHandler(docServ)

		linkRepo = accesslink.NewRepository(db)
		linkServ = accesslink.NewService(linkRepo)
		linkHand = accesslink.NewHandler(linkServ)
	)

	mux.Handle(
		"GET /users/find",
		permissions.RequirePermission(permissions.ReadUsers)(http.HandlerFunc(userHand.Find)),
	)
	mux.Handle(
		"GET /users/auth",
		permissions.RequirePermission(permissions.ReadUsers)(http.HandlerFunc(userHand.Auth)),
	)
	mux.Handle(
		"POST /users",
		permissions.RequirePermission(permissions.CreateUsers)(http.HandlerFunc(userHand.Create)),
	)
	mux.Handle(
		"DELETE /users/{uid}",
		permissions.RequirePermission(permissions.DeleteUsers)(http.HandlerFunc(userHand.Delete)),
	)
	mux.Handle(
		"PATCH /users/{uid}",
		permissions.RequirePermission(permissions.UpdateUsers)(http.HandlerFunc(userHand.Update)),
	)

	mux.Handle(
		"GET /documents",
		permissions.RequirePermission(permissions.ReadDocuments)(http.HandlerFunc(docHand.Get)),
	)
	mux.Handle(
		"GET /documents/{id}",
		permissions.RequirePermission(permissions.ReadDocuments)(http.HandlerFunc(docHand.GetByID)),
	)
	mux.Handle(
		"POST /documents",
		permissions.RequirePermission(permissions.CreateDocuments)(http.HandlerFunc(docHand.Create)),
	)
	mux.Handle(
		"DELETE /documents/{id}",
		permissions.RequirePermission(permissions.DeleteDocuments)(http.HandlerFunc(docHand.Delete)),
	)
	mux.Handle(
		"PATCH /documents/{id}",
		permissions.RequirePermission(permissions.UpdateDocuments)(http.HandlerFunc(docHand.Update)),
	)
	mux.Handle(
		"GET /unsafe/documents/{id}",
		permissions.RequirePermission(permissions.ReadDocumentsPublicly)(http.HandlerFunc(docHand.GetPublicly)),
	)
	mux.Handle(
		"PATCH /unsafe/documents/{id}",
		permissions.RequirePermission(permissions.UpdateDocumentsPublicly)(http.HandlerFunc(docHand.UpdatePublicly)),
	)
	mux.Handle(
		"DELETE /unsafe/documents/{id}",
		permissions.RequirePermission(permissions.DeleteDocumentsPublicly)(http.HandlerFunc(docHand.DeletePublicly)),
	)

	mux.Handle(
		"GET /links",
		permissions.RequirePermission(permissions.ReadLinks)(http.HandlerFunc(linkHand.Get)),
	)
	mux.Handle(
		"GET /links/{token}",
		permissions.RequirePermission(permissions.ReadLinks)(http.HandlerFunc(linkHand.GetByToken)),
	)
	mux.Handle(
		"POST /links",
		permissions.RequirePermission(permissions.CreateLinks)(http.HandlerFunc(linkHand.Create)),
	)
	mux.Handle(
		"DELETE /links/{token}",
		permissions.RequirePermission(permissions.RevokeLinks)(http.HandlerFunc(linkHand.Revoke)),
	)

	// for healthcheck
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
