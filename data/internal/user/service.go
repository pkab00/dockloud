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

func (s *Service) Find(ctx context.Context, req *UserFindRequest) (*UserFindResponse, error) {
	return s.userRepo.Find(ctx, req)
}

func (s *Service) Authenticate(ctx context.Context, req *UserAuthRequest) (*UserResponse, error) {
	var err error

	uid, err := s.userRepo.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Get(ctx, *uid)
	if err != nil {
		return nil, err
	}

	return user.toResponse(), nil
}

func (s *Service) Create(ctx context.Context, req *UserCreateRequest) (*UserResponse, error) {
	if !isValidUsername(req.Username) {
		return nil, apperrors.InavalidNickname
	} else if !isValidEmail(req.Email) {
		return nil, apperrors.InvalidEmail
	}

	u, err := s.userRepo.Create(ctx, req)
	return u.toResponse(), err
}

func (s *Service) Delete(ctx context.Context, uid string) (*UserResponse, error) {
	u, err := s.userRepo.Delete(ctx, uid)
	return u.toResponse(), err
}

func (s *Service) Update(ctx context.Context, uid string, req *UserUpdateRequest) (*UserResponse, error) {
	u, err := s.userRepo.Update(ctx, uid, req)
	return u.toResponse(), err
}
