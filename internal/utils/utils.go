package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if x, ok := data.(error); ok {
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: x.Error()}); err != nil {
			log.Printf("write response failed: %s", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write response failed: %s", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
