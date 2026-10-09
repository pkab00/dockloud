package serviceauth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var serviceKey = os.Getenv("API_SERVICE_KEY")

type ServiceToken struct {
	secret    []byte
	Token     string
	expiresAt time.Time
}

func (st *ServiceToken) Active() bool {
	return time.Now().Before(st.expiresAt)
}

func GenerateToken() (*ServiceToken, error) {
	created := time.Now()
	expires := created.Add(time.Minute * 15)

	claims := jwt.MapClaims{
		"sub": "dockloud-api",            // идентификатор
		"iss": "dockloud-api",            // кем выдан
		"aud": []string{"dockloud-data"}, // кому предназначается
		"iat": created.Unix(),            // когда выдан
		"exp": expires.Unix(),            // когда истекает
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(serviceKey))
	if err != nil {
		return nil, err
	}

	return &ServiceToken{
		secret:    []byte(serviceKey),
		Token:     tokenString,
		expiresAt: expires,
	}, nil
}

func (st *ServiceToken) RefreshIfNeeded() (bool, error) {
	if st.Active() {
		return false, nil
	}

	newToken, err := GenerateToken()
	if err != nil {
		return false, err
	}

	st.Token = newToken.Token
	st.expiresAt = newToken.expiresAt
	return true, nil
}
