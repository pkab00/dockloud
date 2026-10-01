package data

import (
	"database/sql"
	"net/http"

	"github.com/pkab00/dockloud/data/internal/user"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(db *sql.DB, rd *redis.Client, mux *http.ServeMux) {
	var (
		userRepo = user.NewRepositoryWithCache(db, rd)
		userServ = user.NewService(userRepo)
		userHand = user.NewHandler(userServ)
	)

	mux.HandleFunc("GET /users", userHand.GetAll)
	mux.HandleFunc("GET /users/{uid}", userHand.GetByUID)
	mux.HandleFunc("POST /users", userHand.Create)
}
