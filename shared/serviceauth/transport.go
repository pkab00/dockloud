package serviceauth

import (
	"net/http"
)

type AuthTransport struct {
	Base   http.RoundTripper
	Token  *ServiceToken
	Claims Claims
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, err := t.Token.RefreshIfNeeded(t.Claims); err != nil {
		return nil, err
	}

	req = req.Clone(req.Context())
	req.Header.Set("X-Service-Token", t.Token.Token)

	return t.Base.RoundTrip(req)
}
