package document

import "time"

type Response struct {
	ID        int       `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
}

func (d *Document) toResponse() *Response {
	if d == nil {
		return nil
	}

	return &Response{
		ID:        d.ID,
		OwnerID:   d.OwnerID,
		Name:      d.Name,
		Content:   d.Content,
		IsPrivate: d.IsPrivate,
		CreatedAt: d.CreatedAt,
		EditedAt:  d.EditedAt,
	}
}
