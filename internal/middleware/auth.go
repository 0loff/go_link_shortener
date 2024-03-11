// auth модуль проверяет авторизацию пользователя, совершившего обращение к эндпоинту с запросом,
// содержащим cookie Auth и user id, на основании которого происходит аутентификация пользователя.
package middleware

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
	"github.com/0loff/go_link_shortener/pkg/jwt"
	"github.com/google/uuid"
)

// UserAuth - middleware обработчик, проверяющий наличие cookie в запросе.
// После успешного получения Auth cookie из запроса, происходит его разбор и получение user id.
// В случае отсутствия cookie, билдиться новый JWT token и устанавливается в качестве cookie,
// а в service передается предупреждение о запросе от нового пользователя
func UserAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthCookie, err := r.Cookie("Auth")

		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Error("Authentication cookies were not set", zap.Error(err))

				token, err := jwt.BuildJWTString(uuid.New())
				if err != nil {
					logger.Log.Error("Cannot create unique auth token", zap.Error(err))
					panic(err)
				}

				AuthCookie = &http.Cookie{
					Name:  "Auth",
					Value: token,
					Path:  "/",
				}

				http.SetCookie(w, AuthCookie)
				r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUserStatus, "new_user"))

			default:
				logger.Log.Error("Internal server error. Can't get auth cookie from request.", zap.Error(err))
			}
		}

		UserID, err := jwt.GetUserID(AuthCookie.Value)
		if err != nil {
			logger.Log.Error("Failed to get user id from token", zap.Error(err))
		}

		r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUserID, UserID))

		h.ServeHTTP(w, r)
	})
}
