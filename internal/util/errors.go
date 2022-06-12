package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

type JSONErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e JSONErrorMessage) Error() string {
	return e.Message
}

func NewJSONErrorMessage(message, details string) JSONErrorMessage {
	return JSONErrorMessage{
		Message: message,
		Details: details,
	}
}

func JSONError(w http.ResponseWriter, err error, code int) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	var je JSONErrorMessage

	if errors.As(err, &je) {
		json.NewEncoder(w).Encode(je)
	} else {
		json.NewEncoder(w).Encode(JSONErrorMessage{Message: err.Error()})
	}
}
