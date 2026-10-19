package serviceauth

import (
	"errors"
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTTooShortError     = errors.New("JWT secret is too short")
	InvalidTokenError    = errors.New("invalid token")
	MissingSubjectError  = errors.New("missing subject")
	MissingAudienceError = errors.New("missing audience")
)

func containsAudience(provided string, expected []string) bool {
	return slices.Contains(expected, provided)
}

type Claims struct {
	jwt.RegisteredClaims
}

type ServiceConfig struct {
	Secret   []byte
	Subject  string
	Audience string
}

type Verifier struct {
	services map[string]ServiceConfig
}

func NewVerifier(services map[string]ServiceConfig) (*Verifier, error) {
	for _, service := range services {
		if len(service.Secret) < 32 {
			return nil, JWTTooShortError
		}
	}

	return &Verifier{services: services}, nil
}

func (v *Verifier) Verify(rawToken string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		func(t *jwt.Token) (any, error) {
			unverifiedClaims, ok := t.Claims.(*Claims)
			if !ok {
				return nil, InvalidTokenError
			}

			service, ok := v.services[unverifiedClaims.Issuer]
			if !ok {
				return nil, InvalidTokenError
			}

			return service.Secret, nil
		},
		jwt.WithValidMethods(
			[]string{jwt.SigningMethodHS256.Alg()},
		),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", InvalidTokenError
	}

	service := v.services[claims.Issuer]

	if claims.Subject != service.Subject {
		return "", MissingSubjectError
	}
	if !containsAudience(service.Audience, claims.Audience) {
		return "", MissingAudienceError
	}

	return claims.Subject, nil
}
