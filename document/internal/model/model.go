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
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
}

type Permissions struct {
	Read   bool `json:"read"`
	Write  bool `json:"write"`
	Delete bool `json:"delete"`
}

type Link struct {
	DocumentID  int         `json:"document_id"`
	Token       string      `json:"token"`
	Permissions Permissions `json:"permissions"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   *time.Time  `json:"expires_at"`
	RevokedAt   *time.Time  `json:"revoked_at"`
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
