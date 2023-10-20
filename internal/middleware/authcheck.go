package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/h3ll0kitt1/loyality-system/internal/crypto/jwt"
	"github.com/h3ll0kitt1/loyality-system/internal/domain"
	"github.com/h3ll0kitt1/loyality-system/internal/utils"
)

var (
	ErrUnauthorizedAccess = fmt.Errorf("authorized access to this page is forbidden, please, sign-in or sign-up")
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("gohermart")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				utils.WriteResponse(w, http.StatusUnauthorized, ErrUnauthorizedAccess)
			default:
				utils.WriteResponse(w, http.StatusInternalServerError, err)
			}
			return
		}

		var envSecretKey string
		if envSecretKey = os.Getenv("SECRET_KEY"); envSecretKey == "" {
			utils.WriteResponse(w, http.StatusInternalServerError, err)
			return
		}

		login, err := jwt.CheckToken(cookie.Value, envSecretKey)
		if err != nil {
			utils.WriteResponse(w, http.StatusUnauthorized, ErrUnauthorizedAccess)
			return
		}

		ctx := context.WithValue(r.Context(), domain.CtxLoginKey{}, login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
