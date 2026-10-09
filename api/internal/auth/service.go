package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkab00/dockloud/api/internal/client/data"
	apperror "github.com/pkab00/dockloud/api/internal/error"
	"github.com/pkab00/dockloud/api/internal/model"
)

var secretKey = os.Getenv("JWT_KEY")

type ExternalCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ExternalAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Service struct {
	client *data.Client
}

func NewService(client *data.Client) *Service {
	return &Service{client: client}
}

func hashPassword(password string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(password))
	signatureBytes := h.Sum(nil)
	signatureHex := hex.EncodeToString(signatureBytes)

	return signatureHex
}

func generateJWTToken(uid string) (*model.AuthResult, error) {
	created := time.Now()
	expires := created.Add(time.Hour * 2)
	issue := "dockloud-api"

	claims := jwt.MapClaims{
		"sub": uid,            // идентификатор
		"iss": issue,          // кем выдан
		"iat": created.Unix(), // когда выдан
		"exp": expires.Unix(), // когда истекает
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &model.AuthResult{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   time.Until(expires),
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *ExternalCreateRequest) (*model.User, error) {
	var err error

	findRes, err := s.client.FindUser(ctx, &data.FindUserRequest{
		Username: &req.Username,
		Email:    &req.Email,
	})
	if err != nil {
		return nil, err
	}

	switch {
	case findRes.UsernameExists:
		return nil, apperror.UsernameInUse
	case findRes.EmailExists:
		return nil, apperror.EmailInUse
	}

	res, err := s.client.CreateUser(ctx, &data.InternalCreateUserRequest{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashPassword(req.Password),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) AuthUser(ctx context.Context, req *ExternalAuthRequest) (*model.AuthResult, error) {
	var err error

	hashed := hashPassword(req.Password)
	inReq := data.InternalAuthUserRequest{
		Email:        &req.Email,
		PasswordHash: &hashed,
	}
	user, err := s.client.AuthUser(ctx, &inReq)
	if err != nil {
		return nil, err
	}

	token, err := generateJWTToken(user.UID)
	if err != nil {
		return nil, err
	}

	return token, nil
}
