package serviceauth

import (
	"context"
	"net/http"
)

type contextKey struct{}

func Middleware(verifier *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("X-Service-Token")
			res, err := verifier.Verify(header)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				contextKey{},
				res,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ServiceFromContext(ctx context.Context) (string, bool) {
	res, ok := ctx.Value(contextKey{}).(string)
	return res, ok
}
