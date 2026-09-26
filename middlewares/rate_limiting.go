package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

// RateLimitingMiddleware returns per-IP and per-identifier rate limiting for the authentication endpoints. It keys the per-identifier bucket off the
// username (login), clientId (token), or mfaToken (verify-mfa) extracted from the JSON request body. The body is restored so downstream handlers can read it.
//
//goland:noinspection GoUnusedExportedFunction
func RateLimitingMiddleware(requests int, window time.Duration) func(http.Handler) http.Handler {
	perIP := httprate.LimitBy(requests, window, func(r *http.Request) (string, error) {
		ip := middleware.GetClientIP(r.Context())
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		return httprate.CanonicalizeIP(ip), nil
	})

	perIdentifier := httprate.LimitBy(requests, window, authIdentifierKey)

	return func(next http.Handler) http.Handler {
		return perIP(perIdentifier(next))
	}
}

func authIdentifierKey(r *http.Request) (string, error) {
	if r.Body == nil || r.Body == http.NoBody {
		return "", nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	// Restore the body so the handler can decode it again.
	r.Body = io.NopCloser(bytes.NewReader(body))

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil
	}

	var id string
	switch r.URL.Path {
	case "/api/v1/auth/login":
		id = stringValue(payload, "username")
	case "/api/v1/auth/token":
		id = stringValue(payload, "clientId")
	case "/api/v1/auth/verify-mfa":
		id = stringValue(payload, "mfaToken")
	}

	if id == "" {
		return "", nil
	}
	return r.URL.Path + ":" + id, nil
}

func stringValue(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
