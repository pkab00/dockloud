package user

func (u *User) toResponse() *UserResponse {
	return &UserResponse{
		UID:       u.UID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
