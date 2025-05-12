package tokener

import "github.com/omalloc/kratos-admin/pkg/jwt"

type AppToken interface {
	Generate(subjet int64) (string, error)
	Parse(tokenString string) (*jwt.AppClaims, error)
}
