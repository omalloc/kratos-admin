package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
)

type authKey struct{}

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized(reason, "JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized(reason, "keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(reason, "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized(reason, "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized(reason, "Token provider is missing")
	ErrSignToken              = errors.Unauthorized(reason, "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized(reason, "Can not get key while signing token")
)

type AppClaims struct {
	jwt.RegisteredClaims

	UID int64 `json:"uid"`
}

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	signingMethod  jwt.SigningMethod
	tokenHeader    map[string]any
	maxRefreshTime time.Duration // 最大刷新事件
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithTokenHeader withe customer tokenHeader for client side
func WithTokenHeader(header map[string]any) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

// WithMaxRefresh 设置最大刷新时间
func WithMaxRefresh(maxRefreshTime time.Duration) Option {
	return func(o *options) {
		o.maxRefreshTime = maxRefreshTime
	}
}

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(keyFunc jwt.Keyfunc, opts ...Option) middleware.Middleware {
	o := &options{
		signingMethod:  jwt.SigningMethodHS256,
		maxRefreshTime: 7 * 24 * time.Hour, // 7天
	}
	for _, opt := range opts {
		opt(o)
	}

	// got, _ := keyFunc(nil)

	// refreshToken := func(claims *AppClaims) (string, error) {
	// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 	return token.SignedString(got)
	// }

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				if keyFunc == nil {
					return nil, ErrMissingKeyFunc
				}
				auths := strings.SplitN(header.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}
				jwtToken := auths[1]
				var (
					tokenInfo *jwt.Token
					err       error
				)
				tokenInfo, err = jwt.ParseWithClaims(jwtToken, &AppClaims{}, keyFunc)
				if err != nil {
					if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenUnverifiable) {
						return nil, ErrTokenInvalid
					}
					if errors.Is(err, jwt.ErrTokenNotValidYet) {
						return nil, err
					}
					if errors.Is(err, jwt.ErrTokenExpired) {
						// 检查是否过了【最大允许刷新的时间】
						// 首次签名时间 + 最大允许刷新时间区间 > 当前时间 ====> 首次签名时间 > 当前时间 - 最大允许刷新时间区间
						issuedAt, _ := tokenInfo.Claims.GetIssuedAt()
						if issuedAt.Unix() > time.Now().Add(-o.maxRefreshTime).Unix() {
							// 此时并没有过最大允许刷新时间，因此可以重新颁发 token
							// 在这里重新赋值一下过期时间 ExpiresAt 从而达到刷新 token 的目的
							// 但是需要注意的是：一定不能更改 IssuedAt，因为这个字段是用来判断是否过了最大允许刷新时间的
							// refreshToken()
						}

						// 尝试刷新令牌
					}
					return nil, ErrTokenParseFail
				}

				if !tokenInfo.Valid {
					return nil, ErrTokenInvalid
				}
				if tokenInfo.Method != o.signingMethod {
					return nil, ErrUnSupportSigningMethod
				}
				if claims, ok := tokenInfo.Claims.(*AppClaims); ok {
					ctx = NewContext(ctx, claims)
				}
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// Client is a client jwt middleware.
func Client(keyProvider jwt.Keyfunc, opts ...Option) middleware.Middleware {
	claims := &AppClaims{}
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if keyProvider == nil {
				return nil, ErrNeedTokenProvider
			}
			token := jwt.NewWithClaims(o.signingMethod, claims)
			if o.tokenHeader != nil {
				for k, v := range o.tokenHeader {
					token.Header[k] = v
				}
			}
			key, err := keyProvider(token)
			if err != nil {
				return nil, ErrGetKey
			}
			tokenStr, err := token.SignedString(key)
			if err != nil {
				return nil, ErrSignToken
			}
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(authorizationKey, fmt.Sprintf(bearerFormat, tokenStr))
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// NewContext put auth info into context
func NewContext(ctx context.Context, info *AppClaims) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (token *AppClaims, ok bool) {
	token, ok = ctx.Value(authKey{}).(*AppClaims)
	return
}
