package middlewares

import "net/http"

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rr *ResponseRecorder) WriteHeader(code int) {
	rr.StatusCode = code
	rr.ResponseWriter.WriteHeader(code)
}