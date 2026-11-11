package handlers

import (
	"net/http"

	"example.com/go-backend-foundation/pkg/response"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, map[string]string{
		"message": "Hello from Go API",
	})
}
