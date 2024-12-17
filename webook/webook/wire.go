//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
	"webook/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方组件
		ioc.InitRDB, ioc.InitDB, ioc.InitSMSService,

		dao.NewUserDAO,

		cache.NewCodeCache,
		cache.NewUserCache,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		service.NewCodeService,
		service.NewUserService,

		// User 服务
		web.NewUserHandler,

		// 中间件 + server
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
