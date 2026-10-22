package permissions

import (
	"context"
	"net/http"

	"github.com/pkab00/dockloud/shared/contract"
	"github.com/pkab00/dockloud/shared/middleware/service_auth"
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

func checkPermission(ctx context.Context, permission Permission) bool {
	service, ok := serviceauth.ServiceFromContext(ctx)
	if !ok {
		return false
	}
	return permission.CheckFor(service)
}
