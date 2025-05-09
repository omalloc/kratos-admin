package tokener

import "github.com/golang-jwt/jwt/v5"

type AppToken interface {
	Generate(payload map[string]any) (string, error)
	Parse(tokenString string) (jwt.MapClaims, error)
}
