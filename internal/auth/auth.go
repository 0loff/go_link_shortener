package auth

import (
	"context"
	"errors"
	"net/http"

	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/utils"

	"go.uber.org/zap"
)

func UserAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthCookie, err := r.Cookie("Auth")

		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Error("Authentication cookies were not set", zap.Error(err))

				AuthCookie = &http.Cookie{
					Name:  "Auth",
					Value: setAuthToken(),
					Path:  "/",
				}

				http.SetCookie(w, AuthCookie)
				r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUserStatus, "new_user"))

			default:
				logger.Log.Error("Internal server error. Can't get auth cookie from request.", zap.Error(err))
			}
		}

		UserID, err := GetUserID(AuthCookie)
		if err != nil {
			logger.Log.Error("Failed to get user id from token", zap.Error(err))
		}

		r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUserID, UserID))

		h.ServeHTTP(w, r)
	})
}

func setAuthToken() string {
	token, err := BuildJWTString()
	if err != nil {
		logger.Log.Error("Cannot create unique auth token", zap.Error(err))
		panic(err)
	}

	return token
}
