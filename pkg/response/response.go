package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data any `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Error: message,
	})
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, SuccessResponse{
		Data: data,
	})
}
