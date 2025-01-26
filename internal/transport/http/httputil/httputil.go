package httputil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jcserv/go-api-template/internal/utils/log"
)

type HTTPError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func NewHTTPError(code int, message string, details ...map[string]any) *HTTPError {
	detailsMap := map[string]any{}
	for _, detail := range details {
		for k, v := range detail {
			detailsMap[k] = v
		}
	}

	return &HTTPError{
		Code:    code,
		Message: message,
		Details: detailsMap,
	}
}

func BadRequest(w http.ResponseWriter, err error, details ...map[string]any) {
	w.WriteHeader(http.StatusBadRequest)
	writeResponse(w, NewHTTPError(http.StatusBadRequest, err.Error(), details...))
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func InternalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	log.Error(ctx, err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	writeResponse(w, NewHTTPError(http.StatusInternalServerError, err.Error()))
}

func PermanentRedirect(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusPermanentRedirect)
	w.Header().Set("Location", url)
}

func OK(w http.ResponseWriter, response any) {
	w.WriteHeader(http.StatusOK)
	writeResponse(w, response)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeResponse(w http.ResponseWriter, response any) {
	obj, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(obj)
}
