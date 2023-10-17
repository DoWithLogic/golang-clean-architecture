package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils/response"
	"github.com/labstack/echo/v4"
)

// CustomClaims represents the custom claims you want to include in the JWT payload.
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(data CustomClaims, secretKey string) (string, error) {
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthorizeJWT(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, err := extractBearerToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.NewResponseError(http.StatusUnauthorized, response.MsgFailed, err.Error()))
			}

			token, err := jwt.ParseWithClaims(*auth, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.Authentication.Key), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.NewResponseError(http.StatusUnauthorized, response.MsgFailed, err.Error()))
			}

			if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
				c.Set("identity", claims)

				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, response.NewResponseError(http.StatusUnauthorized, response.MsgFailed, err.Error()))
		}
	}
}

func extractBearerToken(c echo.Context) (*string, error) {
	authData := c.Request().Header.Get("Authorization")
	if authData == "" {
		return nil, errors.New("authorization can't be nil")
	}
	parts := strings.Split(authData, " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid authorization value")
	}
	if parts[0] != "Bearer" {
		return nil, errors.New("auth should be bearer")
	}

	return &parts[1], nil
}
