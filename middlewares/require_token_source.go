package middlewares

import (
	"net/http"

	"github.com/floppyman/bs-http/httpres"
)

//goland:noinspection GoUnusedExportedFunction
func RequireTokenSourceMiddleware(source string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uc, ok := GetUserContext(r)
			if !ok {
				httpres.Forbidden(w, "invalid token")
				return
			}

			if uc.TokenSource != source {
				httpres.Forbidden(w, "invalid token source")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
