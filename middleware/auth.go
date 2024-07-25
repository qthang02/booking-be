package middleware

import (
	"github.com/qthang02/booking/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTAuth(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No Authorization header provided"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Could not find bearer token in Authorization header"})
			}

			payload, err := util.ValidateToken(tokenString, secretKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			userInfo, ok := payload.(map[string]interface{})
			if !ok {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid token payload"})
			}

			c.Set("user", userInfo)

			return next(c)
		}
	}
}
