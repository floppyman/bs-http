package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/floppyman/bs-common/jwttoken"
	"github.com/floppyman/bs-common/logging/bsl"

	"github.com/floppyman/bs-http/httpres"
)

//goland:noinspection GoUnusedExportedFunction
func BearerAuthMiddleware(jwtSecret string, jwtIssuer string, getToken func(string) (int, string, error), getUser func(int) (BearerAuthUser, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpres.Unauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httpres.Unauthorized(w, "invalid authorization header format")
				return
			}

			tokenStr := parts[1]

			_, err := jwttoken.VerifyAccessToken(jwtSecret, jwtIssuer, tokenStr)
			if err != nil {
				bsl.Trace("middleware").Err(err).Msg("Invalid JWT token")
				httpres.Unauthorized(w, "invalid token")
				return
			}

			tokenUserId, tokenSource, err := getToken(tokenStr)
			if err != nil {
				bsl.Trace("middleware").Err(err).Msg("Token not found in database")
				httpres.Unauthorized(w, "invalid token")
				return
			}

			user, err := getUser(tokenUserId)
			if err != nil {
				bsl.Trace("middleware").Err(err).Msg("Failed to get user by id")
				httpres.Unauthorized(w, "invalid token")
				return
			}

			uc := UserContext{
				UserId:      user.Id,
				Username:    user.Username,
				Email:       user.Email,
				MfaEnabled:  user.MfaEnabled,
				RoleId:      user.RoleId,
				TokenSource: tokenSource,
			}
			ctx := context.WithValue(r.Context(), UserContextKey{}, uc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type BearerAuthUser struct {
	Id         int
	Username   string
	Email      string
	MfaEnabled bool
	RoleId     int
}
