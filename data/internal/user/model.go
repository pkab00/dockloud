package user

import "time"

type User struct {
	UID          string    `json:"uid" redis:"uid"`
	Username     string    `json:"username" redis:"username"`
	Email        string    `json:"email" redis:"email"`
	PasswordHash string    `json:"password_hash" redis:"password_hash"`
	CreatedAt    time.Time `json:"created_at" redis:"created_at"`
}
