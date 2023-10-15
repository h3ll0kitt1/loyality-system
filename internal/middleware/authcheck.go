package middleware

import (
	"context"
	"errors"
	"net/http"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("gophermarket")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusUnauthorized)
			default:
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		username, ok := checkAuthToken(cookie.Value)
		if !ok {
			http.Error(w, "cookie not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkAuthToken(tokenJWT string) (string, bool) {
	return "", true
}
