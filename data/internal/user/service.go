package user

import (
	"context"
	"net/mail"
	"strings"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

type Service struct {
	userRepo Repository
}

func isValidUsername(username string) bool {
	return len(username) <= 50 &&
		strings.TrimSpace(username) != ""
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

func (s *Service) Find(ctx context.Context, req *UserFindRequest) (*UserFindResult, error) {
	return s.userRepo.Find(ctx, req)
}

func (s *Service) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	if !isValidUsername(req.Username) {
		return nil, apperrors.InavalidNickname
	} else if !isValidEmail(req.Email) {
		return nil, apperrors.InvalidEmail
	}
	return s.userRepo.Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, uid string) (*User, error) {
	return s.userRepo.Delete(ctx, uid)
}

func (s *Service) Update(ctx context.Context, uid string, req *UserUpdateRequest) (*User, error) {
	return s.userRepo.Update(ctx, uid, req)
}
