package model

import "time"

type Document struct {
	ID        int       `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
}

type Link struct {
	Token       string      `json:"token"`
	DocumentID  int         `json:"document_id"`
	CreatorID   string      `json:"creator_id"`
	Permissions Permissions `json:"permissions"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   *time.Time  `json:"expires_at"`
	RevokedAt   *time.Time  `json:"revoked_at"`
}

type Permissions struct {
	Read   *bool `json:"read"`
	Write  *bool `json:"write"`
	Delete *bool `json:"delete"`
}

func (l *Link) Expired() bool {
	now := time.Now()
	return l.ExpiresAt != nil && l.ExpiresAt.Before(now)
}

func (l *Link) Revoked() bool {
	now := time.Now()
	return l.RevokedAt != nil && l.RevokedAt.Before(now)
}
