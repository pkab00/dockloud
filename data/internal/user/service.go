package user

import (
	"context"
	"net/mail"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

type Service struct {
	userRepo Repository
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func NewService(userRepo Repository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) GetAll(ctx context.Context) ([]User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *Service) GetByUID(ctx context.Context, uid string) (*User, error) {
	return s.userRepo.GetByUID(ctx, uid)
}

func (s *Service) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	if !isValidEmail(req.Email) {
		return nil, apperrors.InvalidEmail
	}
	return s.userRepo.Create(ctx, req)
}
