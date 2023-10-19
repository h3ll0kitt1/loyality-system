package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/h3ll0kitt1/loyality-system/internal/crypto/jwt"
)

var (
	ErrUnauthorizedAccess = fmt.Errorf("authorized access to this page is forbidden, please, sign-in or sign-up")
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("gophermarket")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				writeResponse(w, http.StatusUnauthorized, ErrUnauthorizedAccess)
			default:
				writeResponse(w, http.StatusInternalServerError, err)
			}
			return
		}

		login, ok := jwt.CheckToken(cookie.Value)
		if !ok {
			writeResponse(w, http.StatusUnauthorized, ErrUnauthorizedAccess)
			return
		}

		ctx := context.WithValue(r.Context(), "login", login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if x, ok := data.(error); ok {
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: x.Error()}); err != nil {
			log.Printf("write response failed: %w", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write response failed: %w", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
