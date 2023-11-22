package auth

import (
	"go_link_shortener/internal/logger"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const tokenExp = time.Hour * 3
const secretKey = "secretkey"

func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
			},
			UserID: uuid.New(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(authCookie *http.Cookie) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(authCookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Log.Error("The value of the authentication token could not be parsed", zap.Error(err))
		return "", err
	}

	if !token.Valid {
		logger.Log.Error("Invalid auth token param", zap.Error(err))
		return "", err
	}

	return claims.UserID.String(), nil
}
