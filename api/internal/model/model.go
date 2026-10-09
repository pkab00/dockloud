package model

import "time"

type User struct {
	UID       string    `json:"uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Document struct {
	ID        int       `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"is_private"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
}

type UserFindResult struct {
	UIDExists      bool `json:"uid_exists"`
	UsernameExists bool `json:"username_exists"`
	EmailExists    bool `json:"email_exists"`
}

type AuthResult struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   time.Duration `json:"expires_in"`
}

func (r *UserFindResult) UserExists() bool {
	return r.UIDExists || r.UsernameExists || r.EmailExists
}
