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
	Token       string            `json:"token"`
	DocumentID  int               `json:"document_id"`
	CreatorID   string            `json:"creator_id"`
	Permissions AccessPermissions `json:"permissions"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   *time.Time        `json:"expires_at"`
	RevokedAt   *time.Time        `json:"revoked_at"`
}

type AccessPermissions struct {
	Read   bool `json:"read"`
	Write  bool `json:"write"`
	Delete bool `json:"delete"`
}

func (l *Link) isActive() bool {
	now := time.Now()
	if l.ExpiresAt != nil && l.ExpiresAt.Before(now) {
		return false
	} else if l.RevokedAt != nil && l.RevokedAt.Before(now) {
		return false
	}

	return true
}

func (l *Link) isEvergreen() bool {
	return l.ExpiresAt == nil
}
