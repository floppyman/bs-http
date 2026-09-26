package middlewares

import "net/http"

type UserContextKey struct{}

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
	uc, ok := r.Context().Value(UserContextKey{}).(UserContext)
	return uc, ok
}
