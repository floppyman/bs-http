package middlewares

import (
	"net/http"
	"time"

	"github.com/floppyman/bs-common/logging/bsl"
	"github.com/floppyman/bs-common/utils"
)

//goland:noinspection GoUnusedExportedFunction
func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			bsl.
				Trace("request").
				Msgf("%s \033[37m%s\033[0m %s \033[35m%s\033[0m %s",
					utils.ColorMethod(r.Method),
					r.URL.String(),
					utils.ColorStatus(rr.statusCode),
					time.Since(start),
					r.RemoteAddr)
		}()

		next.ServeHTTP(rr, r)
	})
}
