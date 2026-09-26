package middlewares

import (
	"net/http"

	"github.com/floppyman/bs-http/httpres"
)

//goland:noinspection GoUnusedExportedFunction
func RequireRoleIdMiddleware(roleId int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uc, ok := GetUserContext(r)
			if !ok {
				httpres.Forbidden(w, "admin access required")
				return
			}

			if uc.RoleId != roleId {
				httpres.Forbidden(w, "admin access required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
