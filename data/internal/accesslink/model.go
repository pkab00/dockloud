package accesslink

import "time"

type AccessLink struct {
	Token       string            `json:"token" redis:"token"`
	DocumentID  int               `json:"document_id" redis:"document_id"`
	CreatorID   string            `json:"creator_id" redis:"creator_id"`
	Permissions AccessPermissions `json:"permissions" redis:"permissions"`
	CreatedAt   time.Time         `json:"created_at" redis:"created_at"`
	ExpiresAt   *time.Time        `json:"expires_at" redis:"expires_at"`
	RevokedAt   *time.Time        `json:"revoked_at" redis:"revoked_at"`
}

func (al *AccessLink) isActiveLink() bool {
	now := time.Now()
	if al.ExpiresAt != nil && al.ExpiresAt.Before(now) {
		return false
	} else if al.RevokedAt != nil && al.RevokedAt.Before(now) {
		return false
	} else {
		return true
	}
}

type AccessPermissions struct {
	Read   bool `json:"read" redis:"read"`
	Write  bool `json:"write" redis:"write"`
	Delete bool `json:"delete" redis:"delete"`
}
