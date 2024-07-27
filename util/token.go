package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qthang02/booking/enities"
	"time"
)

func GenerateToken(ttl time.Duration, user *enities.User, secretJWTKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})

	tokenString, err := token.SignedString([]byte(secretJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(token string, signedJWTKey string) (map[string]interface{}, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims, nil
}
