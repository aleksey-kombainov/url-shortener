package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

const (
	authCookieName = "auth"
	tokenExpire    = time.Hour * 2 * 24
	secretKey      = "supersecretkey"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type AuthMiddleware struct {
	logger *zerolog.Logger
}

func NewAuthMiddleware(logger *zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{logger: logger}
}

func (m AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {

		var userID string

		cookie, err := request.Cookie(authCookieName)
		if errors.Is(err, http.ErrNoCookie) {
			userID, err = m.generateUserID()
			if err != nil {
				m.logger.Error().Msgf("can't generate userID: %s", err.Error())
			} else {
				jwtStr, err := m.buildJWTString(userID)
				if err != nil {
					m.logger.Error().Msgf("can't generate jwt: %s", err.Error())
				} else {
					cookie := &http.Cookie{
						Name:    authCookieName,
						Value:   jwtStr,
						Path:    "/",
						Expires: time.Now().Add(tokenExpire),
					}
					http.SetCookie(respWriter, cookie)
				}
			}
		} else if err != nil {
			m.logger.Error().Msgf("error while getting cookie: %s", err.Error())
		} else {
			cl, err := m.getClaims(cookie.Value)
			if err != nil {
				m.logger.Error().Msgf("error while getting claims: %s", err.Error())
			} else {
				userID = cl.UserID
			}
		}

		ctx := context.WithValue(request.Context(), user.CtxUserIDKey, userID)
		next.ServeHTTP(respWriter, request.WithContext(ctx))
	})
}

func (m AuthMiddleware) buildJWTString(userID string) (tokenString string, err error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpire)),
		},
		// собственное утверждение
		UserID: userID,
	})
	// создаём строку токена
	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return
	}
	return
}

func (m AuthMiddleware) getClaims(tokenString string) (claims Claims, err error) {
	claims = Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return
	}
	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}
	return
}

func (m AuthMiddleware) generateUserID() (uuidStr string, err error) {
	ud, err := uuid.NewRandom()
	uuidStr = ud.String()
	return
}
