package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/webook/internal/web"
	"webook/webook/internal/web/middleware"
	"webook/webook/pkg/ginx/middlewares/ratelimit"
)

func InitWebServer(mdl []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redis redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl,
		middleware.NewJWTLoginMiddlewareBuilder().
			IgnorePaths("/users/signup").IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms").IgnorePaths("/users/login_sms/code/send").Build(),
		ratelimit.NewBuilder(redis, time.Second, 100).Build(),
	}
}

var corsHdl = func(c *gin.Context) {
	cors.New(cors.Config{
		AllowOrigins: []string{"https://foo.com"},
		AllowMethods: []string{"PUT", "DELETE", "GET", "POST"},
		AllowHeaders: []string{"Content-Type", "authorization"},
		// 允许正式请求中带入的 字段		// 为了 JWT
		ExposeHeaders: []string{"Content-Length", "X-Jwt-Token"},
		// 是否允许 preflight 携带 cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "https://foo.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
