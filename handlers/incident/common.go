package incident

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Message: message,
	})
}
