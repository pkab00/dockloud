package document

import "time"

type Document struct {
	ID        int       `json:"id" redis:"id"`
	OwnerID   string    `json:"owner_id" redis:"owner_id"`
	Name      string    `json:"name" redis:"name"`
	Content   string    `json:"content" redis:"content"`
	IsPrivate bool      `json:"is_private" redis:"is_private"`
	IsDeleted bool      `json:"is_deleted" redis:"is_deleted"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	EditedAt  time.Time `json:"edited_at" redis:"edited_at"`
}
