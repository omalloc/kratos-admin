package tokener

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
func (j *jwtToken) Generate(payload map[string]any) (string, error) {
	if j.opts.secret == "" {
		return "", errors.New("secret is required")
	}

	if (j.opts.ttl / time.Second) <= 0 {
		return "", errors.New("ttl must be greater than 0")
	}

	iat := time.Now().Unix()

	claims := make(jwt.MapClaims)
	claims["exp"] = iat + int64(j.opts.ttl/time.Second)
	claims["iat"] = iat

	j.mergeClaimsPayload(claims, j.opts.payload)
	j.mergeClaimsPayload(claims, payload)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.opts.secret))
}

// Parse implements AppToken.
func (j *jwtToken) Parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.opts.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *jwtToken) mergeClaimsPayload(claims jwt.MapClaims, payload map[string]any) {
	for k, v := range payload {
		claims[k] = v
	}
}
