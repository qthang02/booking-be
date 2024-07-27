package middlewarecustom

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/util"
	"net/http"
	"strings"
)

func JWTAuth(secretKey string, accessibleRoles []string) echo.MiddlewareFunc {
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

			userInfo, err := util.ValidateToken(tokenString, secretKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse user information"})
			}

			userRole := fmt.Sprintf("%v", userInfo["role"])
			if !hasPermission(userRole, accessibleRoles) {
				return c.JSON(http.StatusUnauthorized, "Permission denied")
			}

			c.Set(util.UserID, userInfo["email"])

			return next(c)
		}
	}
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}

	return false
}
