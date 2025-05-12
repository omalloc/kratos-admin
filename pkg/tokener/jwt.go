package tokener

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/omalloc/kratos-admin/pkg/jwt"
)

type jwtOptions struct {
	secret  string         // 密钥
	ttl     time.Duration  // token 存活时长 单位秒
	payload map[string]any // 默认载荷
}

type JwtOption func(*jwtOptions)

type jwtToken struct {
	opts *jwtOptions
}

func WithSecret(secret string) JwtOption {
	return func(o *jwtOptions) {
		o.secret = secret
	}
}

func WithTTL(ttl time.Duration) JwtOption {
	return func(o *jwtOptions) {
		o.ttl = ttl
	}
}

func WithPayload(payload map[string]any) JwtOption {
	return func(o *jwtOptions) {
		o.payload = payload
	}
}

func NewTokener(opts ...JwtOption) AppToken {
	opt := &jwtOptions{
		payload: make(map[string]any),
	}

	for _, apply := range opts {
		apply(opt)
	}

	return &jwtToken{
		opts: opt,
	}
}

// Generate implements AppToken.
func (j *jwtToken) Generate(uid int64) (string, error) {
	if j.opts.secret == "" {
		return "", errors.New("secret is required")
	}

	if (j.opts.ttl / time.Second) <= 0 {
		return "", errors.New("ttl must be greater than 0")
	}

	iat := time.Now()

	claims := &jwt.AppClaims{
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(iat.Add(j.opts.ttl)),
			IssuedAt:  jwtv5.NewNumericDate(iat),
		},
		UID: uid,
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.opts.secret))
}

// Parse implements AppToken.
func (j *jwtToken) Parse(tokenString string) (*jwt.AppClaims, error) {
	token, err := jwtv5.Parse(tokenString, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.opts.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.AppClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *jwtToken) mergeClaimsPayload(claims *jwt.AppClaims, payload map[string]any) {
}
