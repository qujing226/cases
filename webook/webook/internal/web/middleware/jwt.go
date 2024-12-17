package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"webook/webook/internal/web"
)

type JWTLoginMiddlewareBuilder struct {
	paths []string
}

func NewJWTLoginMiddlewareBuilder() *JWTLoginMiddlewareBuilder {
	return &JWTLoginMiddlewareBuilder{}
}

func (j *JWTLoginMiddlewareBuilder) IgnorePaths(path string) *JWTLoginMiddlewareBuilder {
	j.paths = append(j.paths, path)
	return j
}

func (j *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range j.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// authorization 头部
		authCode := ctx.GetHeader("authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authSegments := strings.SplitN(authCode, " ", 2)
		if len(authSegments) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := authSegments[1]
		claims := &web.UserClaims{}
		// ParseWithClaims 里面传指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 验证UserAgent
		if ctx.Request.UserAgent() != claims.UserAgent {
			ctx.String(http.StatusUnauthorized, "attack")
			return
		}

		//  刷新
		if claims.ExpiresAt.Sub(time.Now()) < time.Second*30 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, _ := token.SignedString([]byte("secret"))
			ctx.Header("x-jwt-token", tokenStr)
			return
		}

		ctx.Set("claims", claims)
	}
}
