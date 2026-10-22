package serviceauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ServiceToken struct {
	secret    []byte
	Token     string
	expiresAt time.Time
}

func (st *ServiceToken) Active() bool {
	return time.Now().Before(st.expiresAt)
}

func GenerateToken(secret string, serviceClaims Claims) (*ServiceToken, error) {
	created := time.Now()
	expires := created.Add(time.Minute * 15)

	mapClaims := jwt.MapClaims{
		"sub": serviceClaims.Subject,  // идентификатор
		"iss": serviceClaims.Issuer,   // кем выдан
		"aud": serviceClaims.Audience, // кому предназначается
		"iat": created.Unix(),         // когда выдан
		"exp": expires.Unix(),         // когда истекает
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &ServiceToken{
		secret:    []byte(secret),
		Token:     tokenString,
		expiresAt: expires,
	}, nil
}

func (st *ServiceToken) RefreshIfNeeded(serviceClaims Claims) (bool, error) {
	if st.Active() {
		return false, nil
	}

	newToken, err := GenerateToken(string(st.secret), serviceClaims)
	if err != nil {
		return false, err
	}

	st.Token = newToken.Token
	st.expiresAt = newToken.expiresAt
	return true, nil
}
