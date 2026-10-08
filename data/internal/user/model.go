package user

import "time"

type User struct {
	UID          string    `json:"uid" redis:"uid"`
	Username     string    `json:"username" redis:"username"`
	Email        string    `json:"email" redis:"email"`
	PasswordHash string    `json:"password_hash" redis:"password_hash"`
	CreatedAt    time.Time `json:"created_at" redis:"created_at"`
}

type UserFindResult struct {
	UIDExists      bool `json:"uid_exists"`
	UsernameExists bool `json:"username_exists"`
	EmailExists    bool `json:"email_exists"`
}

func (r *UserFindResult) UserExists() bool {
	return r.UIDExists || r.UsernameExists || r.EmailExists
}
