package middlewares

import (
	"fmt"
	"net/http"

	"github.com/floppyman/bs-common/logging/bsl"
)

// RequestLimitingMiddleware enforces the configured maximum request body size (max receive) and response body size (max send).  A value of 0 disables the corresponding check.
//
//goland:noinspection GoUnusedExportedFunction
func RequestLimitingMiddleware(maxRecv int64, maxSend int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// --- Receive limit --------------------------------------------------
			if maxRecv > 0 {
				r.Body = http.MaxBytesReader(w, r.Body, maxRecv)
			}

			// --- Send limit -----------------------------------------------------
			var wrapped *limitResponseWriter
			if maxSend > 0 {
				wrapped = &limitResponseWriter{
					ResponseWriter: w,
					limit:          maxSend,
					reached:        false,
				}
				w = wrapped
			}

			next.ServeHTTP(w, r)

			// If the response limit was hit, log an error.
			if wrapped != nil && wrapped.reached {
				bsl.Errorf("request-limits", "Response body exceeded configured max_send_bytes (limit=%d) on %s", maxSend, r.URL.Path)
			}
		})
	}
}

// limitResponseWriter wraps an http.ResponseWriter and counts the bytes
// written.  When the configured limit is exceeded it stops writing and
// records the condition.
type limitResponseWriter struct {
	http.ResponseWriter
	limit   int64
	written int64
	reached bool
}

func (l *limitResponseWriter) Write(p []byte) (int, error) {
	if l.reached {
		return 0, fmt.Errorf("response body exceeded configured max_send_bytes limit (%d bytes)", l.limit)
	}

	l.written += int64(len(p))
	if l.written > l.limit {
		l.reached = true
		// The headers may already be flushed, so we cannot change the
		// status code.  Just stop writing and return an error.
		return 0, fmt.Errorf("response body exceeded configured max_send_bytes limit (%d bytes)", l.limit)
	}
	return l.ResponseWriter.Write(p)
}
