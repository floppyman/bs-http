package httpres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func sendError(w http.ResponseWriter, code int, res string) {
	http.Error(w, fmt.Sprintf(`{"error":"%s"}`, res), code)
}

//goland:noinspection GoUnusedExportedFunction
func Unauthorized(w http.ResponseWriter, res string) {
	sendError(w, http.StatusUnauthorized, res)
}

//goland:noinspection GoUnusedExportedFunction
func Forbidden(w http.ResponseWriter, res string) {
	sendError(w, http.StatusForbidden, res)
}

//goland:noinspection GoUnusedExportedFunction
func NotAcceptable(w http.ResponseWriter, allowedMethods []string) error {
	w.Header().Add("Allow", strings.Join(allowedMethods, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte(""))
	if err != nil {
		return err
	}
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func BadRequestJson(w http.ResponseWriter, res any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func NotFoundJson(w http.ResponseWriter, res any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func BadRequestString(w http.ResponseWriter, res string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func OkJson(w http.ResponseWriter, res any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

// ServiceUnavailableJson writes res as JSON with HTTP 503.
//
//goland:noinspection GoUnusedExportedFunction
func ServiceUnavailableJson(w http.ResponseWriter, res any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func OkString(w http.ResponseWriter, res string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}
