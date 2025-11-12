package middleware

import (
	"context"
	"sync"
	"time"
	"west2/pkg/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

var (
	jwtMiddleware *jwt.HertzJWTMiddleware
	jwtOnce       sync.Once
	initErr       error
	identityKey   string
)

func GetJWTMiddleware() (*jwt.HertzJWTMiddleware, error) {
	jwtOnce.Do(func() {
		cfg := config.GetConfig()
		identityKey = "uid"
		jwtMiddleware, initErr = jwt.New(&jwt.HertzJWTMiddleware{
			Key:           []byte(cfg.Jwt.SecretKey),
			Timeout:       time.Hour * cfg.Jwt.AccessTimeout,
			MaxRefresh:    time.Hour * cfg.Jwt.RefreshTimeout,
			TokenLookup:   "header:Access-Token",
			TokenHeadName: "Bearer",
			IdentityKey:   identityKey,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(string); ok {
					return jwt.MapClaims{identityKey: v}
				}
				return jwt.MapClaims{}
			},
			IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
				claims := jwt.ExtractClaims(ctx, c)
				return claims[identityKey].(string)
			},
		})

	})
	return jwtMiddleware, initErr
}

func GenerateToken(uid string) (string, time.Time, error) {
	middleware, err := GetJWTMiddleware()
	if err != nil {
		return "", time.Time{}, err
	}
	return middleware.TokenGenerator(uid)
}

func GetUserFromContext(ctx context.Context, c *app.RequestContext) string {
	claims := jwt.ExtractClaims(ctx, c)

	if uid, ok := claims[identityKey].(string); ok {
		return uid
	}
	return ""
}
