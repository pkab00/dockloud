package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/pkab00/dockloud/api/internal/client/data"
	"github.com/pkab00/dockloud/api/internal/model"
)

type ExternalCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Service struct {
	userClient *data.Client
}

func NewService(client *data.Client) *Service {
	return &Service{userClient: client}
}

func hashPassword(password string) string {
	key := os.Getenv("JWT_KEY")

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(password))
	signatureBytes := h.Sum(nil)
	signatureHex := hex.EncodeToString(signatureBytes)

	return signatureHex
}

func (s *Service) CreateUser(ctx context.Context, req *ExternalCreateRequest) (*model.User, error) {
	internalReq := data.InternalCreateRequest{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashPassword(req.Password),
	}

	res, err := s.userClient.CreateUser(ctx, &internalReq)
	if err != nil {
		return nil, err
	}

	return res, nil
}
