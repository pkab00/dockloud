package serviceauth

import (
	"net/http"

	"github.com/pkab00/dockloud/shared/contract"
)

func RequirePermission(permission Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !checkPermission(r.Context(), permission) {
				contract.WriteErrorAsJSON(
					w,
					http.StatusForbidden,
					contract.ErrPermissionDenied,
					"The service has no permission to access this resource",
				)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
