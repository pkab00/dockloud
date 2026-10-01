package user

import "context"

type Service struct {
	userRepo Repository
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
