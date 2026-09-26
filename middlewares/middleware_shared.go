package middlewares

import "net/http"

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

type userContextKey struct{}

// UserContext holds the authenticated user's information stored in the request context.
type UserContext struct {
	UserId      int
	Username    string
	Email       string
	MfaEnabled  bool
	RoleId      int
	TokenSource string
}

// GetUserContext extracts the UserContext from the request context.
func GetUserContext(r *http.Request) (UserContext, bool) {
	uc, ok := r.Context().Value(userContextKey{}).(UserContext)
	return uc, ok
}
