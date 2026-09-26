package middlewares

import (
	"net/http"
	"slices"

	"github.com/floppyman/bs-http/httpres"
)

//goland:noinspection GoUnusedExportedFunction
func MethodsNotAcceptableMiddleware(methods []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(methods, r.Method) {
				next.ServeHTTP(w, r)
				return
			}
			_ = httpres.NotAcceptable(w, methods)
		})
	}
}
