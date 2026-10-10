package user

import "time"

type UserResponse struct {
	UID       string    `json:"uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserFindResponse struct {
	UIDExists      bool `json:"uid_exists"`
	UsernameExists bool `json:"username_exists"`
	EmailExists    bool `json:"email_exists"`
}

func (r *UserFindResponse) UserExists() bool {
	return r.UIDExists || r.UsernameExists || r.EmailExists
}

func (u *User) toResponse() *UserResponse {
	return &UserResponse{
		UID:       u.UID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
